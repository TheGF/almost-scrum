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

function setPendingTasks() {
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

const errorHandlers = []

function errorHandler(r) {
    for (const [_,handler] of errorHandlers) {
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


    static createTask(project, board, title) {
        title = encodeURIComponent(title)
        return axios.post(`/api/v1/projects/${project}/boards/${board}?title=${title}`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static setTaskLater(project, board, name, content) {
        return new Promise(function (resolve, reject) {
            const k = `setTask:${project}/${board}/${name}`
            pendingSet[k] = [Date.now() + setDelay, [resolve, reject], Server.setTask, project, board, name, content]
            if (pendingInterval == null) {
                pendingInterval = setInterval(setPendingTasks, setDelay)
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
            const k = `uploadFile:${project}/${path}/${name}`
            pendingSet[k] = [Date.now() + setDelay, [resolve, reject], Server.uploadFileToLibrary, project,
                path, file, name]
            if (pendingInterval == null) {
                pendingInterval = setInterval(setPendingTasks, setDelay)
            }
        })
    }

    static deleteFromLibrary(project, path, recursive) {
        path = encodeURIComponent(path)
        let url = `/api/v1/projects/${project}/library${path}`
        if (recursive) { url += '?recursive' }
        return axios.delete(url, getConfig(url))
            .then(r => r.data)
            .catch(errorHandler);
    }

    static downloadFromlibrary(project, path) {
        path = encodeURIComponent(path)
        return axios.get(`/api/v1/projects/${project}/library${path}`, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static openFromlibrary(project, path) {
        path = encodeURIComponent(path)

        const link = document.createElement("a");
        link.href = `/api/v1/projects/${project}/library${path}?token=${localStorage.token}`;
        link.target = '_';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    static listLibrary(project, path) {
        path = encodeURIComponent(path)
        return axios.get(`/api/v1/projects/${project}/library${path}`, getConfig())
            .then(r => r.data.filter(f => !f.name.startsWith('.')))
            .catch(errorHandler);
    }

    static createFolderInLibrary(project, path) {
        path = encodeURIComponent(path)
        return axios.put(`/api/v1/projects/${project}/library${path}`, null, getConfig())
            .then(r => r.data)
            .catch(errorHandler);
    }

    static moveFileInLibrary(project, oldpath, path) {
        oldpath = encodeURIComponent(oldpath)
        path = encodeURIComponent(path)
        return axios.post(`/api/v1/projects/${project}/library${path}?move=${oldpath}`,
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

}

export default Server;