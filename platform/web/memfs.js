// In-memory filesystem for Tyumi's browser build.
// Load this script before wasm_exec.js. Go's js/wasm port calls globalThis.fs for
// os.Open, os.WriteFile, and directory listings, and the stock wasm_exec.js
// filesystem does not implement those calls.
(function () {
  const files = new Map();
  files.set("/", { dir: true, data: new Uint8Array(0), mtime: Date.now() });

  let cwd = "/";
  let nextFd = 3;
  const fds = new Map();

  function norm(p) {
    const parts = [];
    for (const part of String(p).split("/")) {
      if (part === "" || part === ".") continue;
      if (part === "..") parts.pop();
      else parts.push(part);
    }
    return "/" + parts.join("/");
  }

  function resolve(p) {
    if (!p) return cwd;
    if (p.startsWith("/")) return norm(p);
    return norm(cwd + "/" + p);
  }

  function error(code) {
    const err = new Error(code);
    err.code = code;
    return err;
  }

  function parent(path) {
    const n = norm(path);
    if (n === "/") return "/";
    return n.slice(0, n.lastIndexOf("/")) || "/";
  }

  function base(path) {
    const n = norm(path);
    return n.slice(n.lastIndexOf("/") + 1);
  }

  function statOf(node) {
    const mode = (node.dir ? 0x4000 : 0x8000) | 0o755;
    return {
      dev: 1,
      ino: 1,
      mode: mode,
      nlink: 1,
      uid: 0,
      gid: 0,
      rdev: 0,
      size: node.data.length,
      blksize: 4096,
      blocks: Math.ceil(node.data.length / 512),
      atimeMs: node.mtime,
      mtimeMs: node.mtime,
      ctimeMs: node.mtime,
      isDirectory: function () { return node.dir; },
    };
  }

  function children(dir) {
    const prefix = dir === "/" ? "/" : dir + "/";
    const names = [];
    for (const key of files.keys()) {
      if (key === dir || !key.startsWith(prefix)) continue;
      const rest = key.slice(prefix.length);
      if (!rest.includes("/")) names.push(rest);
    }
    return names;
  }

  const constants = {
    O_RDONLY: 0,
    O_WRONLY: 1,
    O_RDWR: 2,
    O_CREAT: 64,
    O_EXCL: 128,
    O_TRUNC: 512,
    O_APPEND: 1024,
    O_DIRECTORY: 65536,
  };

  let outputBuf = "";
  const decoder = new TextDecoder("utf-8");

  function writeSync(fd, buf) {
    if (fd === 1 || fd === 2) {
      outputBuf += decoder.decode(buf, { stream: true });
      let nl = outputBuf.indexOf("\n");
      while (nl !== -1) {
        const line = outputBuf.slice(0, nl);
        outputBuf = outputBuf.slice(nl + 1);
        if (fd === 2) console.error(line);
        else console.log(line);
        nl = outputBuf.indexOf("\n");
      }
      return buf.length;
    }
    const f = fds.get(fd);
    if (!f || f.node.dir) return 0;
    const bytes = buf instanceof Uint8Array ? buf : new Uint8Array(buf);
    return writeBytes(f, bytes, 0, bytes.length, null);
  }

  function writeBytes(f, buf, offset, length, position) {
    const chunk = buf.subarray(offset, offset + length);
    let start = position == null ? f.pos : position;
    if (position == null && (f.flags & constants.O_APPEND)) start = f.node.data.length;
    const end = start + chunk.length;
    if (f.node.data.length < end) {
      const next = new Uint8Array(end);
      next.set(f.node.data);
      f.node.data = next;
    }
    f.node.data.set(chunk, start);
    if (position == null) f.pos = start + chunk.length;
    f.node.mtime = Date.now();
    return chunk.length;
  }

  const fs = {
    constants: constants,
    writeSync: writeSync,
    write: function (fd, buf, offset, length, position, cb) {
      try {
        const f = fds.get(fd);
        if (!f) return cb(error("EBADF"));
        if (f.node.dir) return cb(error("EISDIR"));
        cb(null, writeBytes(f, buf, offset, length, position));
      } catch (err) {
        cb(err);
      }
    },
    read: function (fd, buffer, offset, length, position, cb) {
      const f = fds.get(fd);
      if (!f) return cb(error("EBADF"));
      if (f.node.dir) return cb(error("EISDIR"));
      const start = position == null ? f.pos : position;
      const n = Math.max(0, Math.min(length, f.node.data.length - start));
      if (n > 0) buffer.set(f.node.data.subarray(start, start + n), offset);
      if (position == null) f.pos = start + n;
      cb(null, n);
    },
    open: function (path, flags, mode, cb) {
      const p = resolve(path);
      let node = files.get(p);
      const creat = (flags & constants.O_CREAT) !== 0;
      const excl = (flags & constants.O_EXCL) !== 0;
      const trunc = (flags & constants.O_TRUNC) !== 0;
      const dirOnly = (flags & constants.O_DIRECTORY) !== 0;
      if (!node) {
        if (!creat) return cb(error("ENOENT"));
        if (!files.get(parent(p)) || !files.get(parent(p)).dir) return cb(error("ENOENT"));
        node = { dir: false, data: new Uint8Array(0), mtime: Date.now() };
        files.set(p, node);
      } else if (creat && excl) {
        return cb(error("EEXIST"));
      }
      if (dirOnly && !node.dir) return cb(error("ENOTDIR"));
      if (trunc && !node.dir) node.data = new Uint8Array(0);
      const fd = nextFd++;
      fds.set(fd, { node: node, path: p, pos: 0, flags: flags });
      cb(null, fd);
    },
    close: function (fd, cb) {
      fds.delete(fd);
      cb(null);
    },
    fstat: function (fd, cb) {
      const f = fds.get(fd);
      if (!f) return cb(error("EBADF"));
      cb(null, statOf(f.node));
    },
    stat: function (path, cb) {
      const node = files.get(resolve(path));
      if (!node) return cb(error("ENOENT"));
      cb(null, statOf(node));
    },
    lstat: function (path, cb) { fs.stat(path, cb); },
    mkdir: function (path, perm, cb) {
      const p = resolve(path);
      if (files.has(p)) return cb(error("EEXIST"));
      const dir = files.get(parent(p));
      if (!dir || !dir.dir) return cb(error("ENOENT"));
      files.set(p, { dir: true, data: new Uint8Array(0), mtime: Date.now() });
      cb(null);
    },
    readdir: function (path, cb) {
      const p = resolve(path);
      const node = files.get(p);
      if (!node) return cb(error("ENOENT"));
      if (!node.dir) return cb(error("ENOTDIR"));
      cb(null, children(p));
    },
    rmdir: function (path, cb) {
      const p = resolve(path);
      const node = files.get(p);
      if (!node) return cb(error("ENOENT"));
      if (!node.dir) return cb(error("ENOTDIR"));
      if (children(p).length > 0) return cb(error("ENOTEMPTY"));
      if (p === "/") return cb(error("EPERM"));
      files.delete(p);
      cb(null);
    },
    unlink: function (path, cb) {
      const p = resolve(path);
      const node = files.get(p);
      if (!node) return cb(error("ENOENT"));
      if (node.dir) return cb(error("EISDIR"));
      files.delete(p);
      cb(null);
    },
    rename: function (from, to, cb) {
      const a = resolve(from);
      const b = resolve(to);
      const node = files.get(a);
      if (!node) return cb(error("ENOENT"));
      if (!files.get(parent(b)) || !files.get(parent(b)).dir) return cb(error("ENOENT"));
      files.delete(a);
      files.set(b, node);
      cb(null);
    },
    fsync: function (fd, cb) { cb(null); },
    ftruncate: function (fd, length, cb) {
      const f = fds.get(fd);
      if (!f) return cb(error("EBADF"));
      const next = new Uint8Array(length);
      next.set(f.node.data.subarray(0, length));
      f.node.data = next;
      cb(null);
    },
    truncate: function (path, length, cb) {
      const node = files.get(resolve(path));
      if (!node) return cb(error("ENOENT"));
      const next = new Uint8Array(length);
      next.set(node.data.subarray(0, length));
      node.data = next;
      cb(null);
    },
    chmod: function (path, mode, cb) { cb(null); },
    fchmod: function (fd, mode, cb) { cb(null); },
    chown: function (path, uid, gid, cb) { cb(null); },
    fchown: function (fd, uid, gid, cb) { cb(null); },
    utimes: function (path, atime, mtime, cb) {
      const node = files.get(resolve(path));
      if (node) node.mtime = Date.now();
      cb(null);
    },
  };

  globalThis.fs = fs;
  globalThis.process = {
    cwd: function () { return cwd; },
    chdir: function (p) { cwd = resolve(p); },
    getuid: function () { return 0; },
    getgid: function () { return 0; },
    geteuid: function () { return 0; },
    getegid: function () { return 0; },
    getgroups: function () { return []; },
    pid: 1,
    ppid: 0,
    umask: function () { return 0; },
  };
  globalThis.path = {
    resolve: function () {
      let out = "";
      for (let i = 0; i < arguments.length; i++) {
        const seg = arguments[i];
        if (!seg) continue;
        if (seg.startsWith("/")) out = seg;
        else out = (out || cwd) + "/" + seg;
      }
      return norm(out || cwd);
    },
  };
})();
