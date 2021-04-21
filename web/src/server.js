import axios from 'axios';


function getConfig() {
    const token = localStorage.token || '';

    return token ? {
        headers: {
            Authorization: `Bearer ${token}`,
        }
    } : {}
}

const pendingSet = {}
const setDelay = 2 * 1000
let pendingInterval = null


const errorHandlers = []

function errorHandler(r) {
    for (const [_, handler] of errorHandlers) {
        r = handler(r)
        if (!r) return
    }

    return Promise.reject(r)
}

class Server {

    static addErrorHandler(priority, handler) {
        if (errorHandlers.filter(i => i[1] == handler).length)
            return
        errorHandlers.push([priority, handler])
        errorHandlers.sort((a, b) => b[0] - a[0])
    }

    static sendPendingTasks() {
        if (pendingSet.length == 0) {
            const interval = pendingInterval
            pendingInterval = null
            cancelInterval(interval)
        }
        for (const k in pendingSet) {
            const [tm, [resolve, reject], f, ...params] = pendingSet[k];
            if (tm < Date.now()) {
                delete pendingSet[k];

                f(...params)
                    .then(resolve)
                    .catch(reject);
            }
        }
    }


    static authenticate(username, password) {
        const credentials = {
            username: username,
            password: password,
        }
        return axios.post("/auth/login", credentials, getConfig())
            .then(r => {
                const token = r.data && r.data.token;
                localStorage.setItem("token", token)
                return token
            })
    }

    static getResource(path) {
        return axios.get(path, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static hello(id) {
        return axios.post(`/auth/hello?id=${id}`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static bye(id) {
        return axios.post(`/auth/bye?id=${id}`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }


    static getLocalUsers() {
        return axios.get('/api/v1/passwords', getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postLocalUserCredentials(username, password) {
        const credentials = {
            username: username,
            password: password,
        }
        return axios.post('/api/v1/passwords', credentials, getConfig())
            .then(r => r.data)
            .catch(errorHandler)
    }

    static getLoginUser() {
        return axios.get('/api/v1/user', getConfig())
            .then(r => r.data)
            .catch(errorHandler)
    }

    static getProjectsList() {
        return axios.get('/api/v1/projects', getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getTemplatesList() {
        return axios.get('/api/v1/templates', getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static createProject(name, templates) {
        const body = {
            projectName: name,
            templates: templates,
        }
        return axios.post(`/api/v1/projects`, body, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static importProject(path, inject, templates) {
        const body = {
            importPath: path,
            inject: inject,
            templates: templates,
        }
        return axios.post(`/api/v1/projects`, body, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static cloneFromGit(url, inject, templates) {
        const body = {
            gitUrl: url,
            inject: inject,
            templates: templates,
        }
        return axios.post(`/api/v1/projects`, body, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }



    static getProjectInfo(project) {
        return axios.get(`/api/v1/projects/${project}/info`, getConfig())
            .then(r => r.data)
            .catch(errorHandler)
    }

    static setProjectInfo(project, info) {
        return axios.put(`/api/v1/projects/${project}/info`, info, getConfig())
            .then(r => r.data)
            .catch(errorHandler)
    }

    static listUsers(project) {
        return axios.get(`/api/v1/projects/${project}/users`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static setUser(project, user, userInfo) {
        return axios.put(`/api/v1/projects/${project}/users/${user}`,
            userInfo, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getUser(project, user) {
        return axios.get(`/api/v1/projects/${project}/users/${user}`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static delUser(project, user) {
        return axios.delete(`/api/v1/projects/${project}/users/${user}`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static listBoards(project) {
        return axios.get(`/api/v1/projects/${project}/boards`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static createBoard(project, board) {
        return axios.put(`/api/v1/projects/${project}/boards/${board}`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static renameBoard(project, oldName, newName) {
        return axios.put(`/api/v1/projects/${project}/boards/${newName}?rename=${oldName}`,
            null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static deleteBoard(project, board) {
        return axios.delete(`/api/v1/projects/${project}/boards/${board}`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }


    static listTasks(project, board, filter, start, end) {
        let url = `/api/v1/projects/${project}/boards/${board}?`
        let params = []
        if (start) params.push(`start=${start}`);
        if (end) params.push(`end=${end}`);
        if (filter) params.push(`filter=${encodeURIComponent(filter)}`);
        url += params.join('&')

        return axios.get(url, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getTask(project, board, name) {
        name = encodeURIComponent(name)
        return axios.get(`/api/v1/projects/${project}/boards/${board}/${name}`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }


    static createTask(project, board, title, type) {
        title = encodeURIComponent(title)
        return axios.post(`/api/v1/projects/${project}/boards/${board}?title=${title}&type=${type}`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static setTaskLater(project, board, name, content) {
        return new Promise(function (resolve, reject) {
            const k = `setTask:${project}/${board}/${name}`
            pendingSet[k] = [Date.now() + setDelay, [resolve, reject], Server.setTask, project, board, name, content]
            if (pendingInterval == null) {
                pendingInterval = setInterval(Server.sendPendingTasks, setDelay)
            }
        })
    }

    static setTask(project, board, name, content) {
        name = encodeURIComponent(name)
        return axios.put(`/api/v1/projects/${project}/boards/${board}/${name}`, content, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static moveTask(project, board, name, newBoard, title) {
        name = encodeURIComponent(name)
        let url = `/api/v1/projects/${project}/boards/${newBoard}?move=${board}/${name}`
        if (title) {
            url += `&title=${encodeURIComponent(title)}`;
        }
        return axios.post(url,
            null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static deleteTask(project, board, name) {
        name = encodeURIComponent(name)
        return axios.delete(`/api/v1/projects/${project}/boards/${board}/${name}`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }



    static touchTask(project, board, name) {
        name = encodeURIComponent(name)
        return axios.post(`/api/v1/projects/${project}/boards/${board}/${name}?touch`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler)
    }

    static uploadFileToLibrary(project, path, file, name) {
        const token = localStorage.token
        path = encodeURIComponent(path)
        const config = getConfig();
        let formData = null;

        if (file) {
            formData = new FormData();
            if (name) {
                formData.append("file", file, name);
            } else {
                formData.append("file", file);
            }
            config.headers = config.headers || {};
            config.headers['Content-Type'] = 'multipart/form-data';
        }

        return axios.post(`/api/v1/projects/${project}/library${path}?token=${token}`,
            formData, { config })
            .then(r => r.data.filter(f => !f.name.startsWith('.')))
            .catch(errorHandler);
    }

    static uploadFileToLibraryLater(project, path, file, name) {
        return new Promise(function (resolve, reject) {
            const k = name ? `uploadFile:${project}/${path}/${name}` :
                `uploadFile:${project}/${path}`
            pendingSet[k] = [Date.now() + setDelay, [resolve, reject], Server.uploadFileToLibrary, project,
                path, file, name
            ]
            if (pendingInterval == null) {
                pendingInterval = setInterval(Server.sendPendingTasks, setDelay)
            }
        })
    }

    static deleteFromLibrary(project, path, recursive = false, archive = false) {
        path = encodeURIComponent(path)
        const target = archive ? 'archive' : 'library'
        let url = `/api/v1/projects/${project}/${target}${path}`
        if (recursive) { url += '?recursive' }
        return axios.delete(url, getConfig(url))
            .then(r => r.data)
            .catch(errorHandler);
    }

    static downloadFromlibrary(project, path, archive = false) {
        path = encodeURIComponent(path)
        const target = archive ? 'archive' : 'library'
        return axios.get(`/api/v1/projects/${project}/${target}${path}`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static localOpenFromLibrary(project, path, archive = false) {
        path = encodeURIComponent(path)
        const target = archive ? 'archive' : 'library'
        return axios.get(`/api/v1/projects/${project}/${target}${path}?local`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }


    static openFromlibrary(project, path, archive = false) {
        path = encodeURIComponent(path)
        const target = archive ? 'archive' : 'library'
        const link = document.createElement("a");
        link.href = `/api/v1/projects/${project}/${target}${path}?token=${localStorage.token}`;
        link.target = '_';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    static postNewBook(project, path) {
        const link = document.createElement("a");
        link.href = `/api/v1/projects/${project}/library-book${path}?token=${localStorage.token}`;
        link.target = '_';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    static postNewBook2(project, path, settings) {
        return axios.post(`/api/v1/projects/${project}/library-book${path}`, settings, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }


    static listLibrary(project, path, archive = false) {
        path = encodeURIComponent(path)
        const target = archive ? 'archive' : 'library'
        return axios.get(`/api/v1/projects/${project}/${target}${path}`, getConfig())
            .then(r => r.data.filter(f => !f.name.startsWith('.')))
            .catch(errorHandler);
    }

    static getVersions(project, path) {
        path = encodeURIComponent(path)
        return axios.get(`/api/v1/projects/${project}/library${path}?versions`, getConfig())
            .then(r => r.data.filter(f => !f.name.startsWith('.')))
            .catch(errorHandler);
    }

    static createFolderInLibrary(project, path, archive = false) {
        path = encodeURIComponent(path)
        const target = archive ? 'archive' : 'library'
        return axios.put(`/api/v1/projects/${project}/${target}${path}`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static moveFileInLibrary(project, oldpath, path) {
        oldpath = encodeURIComponent(oldpath)
        path = encodeURIComponent(path)
        return axios.post(`/api/v1/projects/${project}/library${path}?action=move&origin=${oldpath}`,
            null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static upgradeVersion(project, path) {
        path = encodeURIComponent(path)
        return axios.post(`/api/v1/projects/${project}/library${path}?action=upgrade`,
            null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static setVisibility(project, path, public_) {
        path = encodeURIComponent(path)
        return axios.post(`/api/v1/projects/${project}/library${path}?action=visibility&public=${public_}`,
            null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getLibraryStat(project, files) {
        return axios.post(`/api/v1/projects/${project}/library-stat`,
            files, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }


    static getSuggestions(project, prefix, total) {
        prefix = encodeURIComponent(prefix)
        let url = `/api/v1/projects/${project}/index/suggest/${prefix}`
        if (total) url += `&total=${total}`
        return axios.get(url, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getFedConfig(project) {
        let url = `/api/v1/projects/${project}/fed/config`
        return axios.get(url, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postFedConfig(project, config) {
        let url = `/api/v1/projects/${project}/fed/config`
        return axios.post(url, config, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postFedShare(project, key, exchanges, removeCredentials) {
        const c = {
            key: key,
            exchanges: exchanges,
            removeCredentials: removeCredentials,
        }
        let url = `/api/v1/projects/${project}/fed/share`
        return axios.post(url, c, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postFedJoin(project, key, token) {
        const invite = {
            key: key,
            token: token,
        }
        let url = `/api/v1/projects/${project}/fed/join`
        return axios.post(url, invite, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }


    static getFedStatus(project) {
        let url = `/api/v1/projects/${project}/fed/status`
        return axios.get(url, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getFedDiffs(project, sync) {
        let url = `/api/v1/projects/${project}/fed/diffs`
        if (sync) url += '?sync'
        return axios.get(url, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postFedImport(project, logs) {
        let url = `/api/v1/projects/${project}/fed/import`
        return axios.post(url, logs, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postFedExport(project, since = null) {
        let url = `/api/v1/projects/${project}/fed/export`
        if (since) url += `?since=${since.toISOString()}`
        return axios.post(url, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postFedPull(project) {
        let url = `/api/v1/projects/${project}/fed/pull`
        return axios.post(url, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postFedPush(project) {
        let url = `/api/v1/projects/${project}/fed/push`
        return axios.post(url, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getGitStatus(project) {
        return axios.get(`/api/v1/projects/${project}/git/status`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postGitCommit(project, commitInfo) {
        return axios.post(`/api/v1/projects/${project}/git/commit`, commitInfo, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postGitPush(project) {
        return axios.post(`/api/v1/projects/${project}/git/push`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static postGitPull(project) {
        return axios.post(`/api/v1/projects/${project}/git/pull`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getGitSettings(project) {
        return axios.get(`/api/v1/projects/${project}/git/settings`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static putGitSettings(project, settings) {
        return axios.put(`/api/v1/projects/${project}/git/settings`, settings, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static getGanttTasks(project) {
        return axios.get(`/api/v1/projects/${project}/gantt`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

}

export default Server;