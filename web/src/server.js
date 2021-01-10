import axios from 'axios';

function loginWhenUnauthorized(r) {
    if (r && r.response && r.response.status == 401) {
        localStorage.removeItem('username');
        localStorage.removeItem('token');
        window.location.assign(window.location.href);
    }
}

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
        const [tm, project, board, name, content] = pendingSet[k];
        if (tm < Date.now()) {
            delete pendingSet[k];
            Server.setTask(project, board, name, content);
        }
    }
}

class Server {

    static getProjectsList() {
        return axios.get('/api/v1/projects', getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static getProjectInfo(project) {
        return axios.get(`/api/v1/projects/${project}/info`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static listUsers(project) {
        return axios.get(`/api/v1/projects/${project}/users`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static listBoards(project) {
        return axios.get(`/api/v1/projects/${project}/boards`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static createBoard(project, board) {
        return axios.put(`/api/v1/projects/${project}/boards/${board}`, null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
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
            .catch(loginWhenUnauthorized);
    }

    static getTask(project, board, name) {
        name = encodeURIComponent(name)
        return axios.get(`/api/v1/projects/${project}/boards/${board}/${name}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }


    static createTask(project, board, title) {
        title = encodeURIComponent(title)
        return axios.post(`/api/v1/projects/${project}/boards/${board}?title=${title}`, null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static setTaskLater(project, board, name, content) {
        const k = `${project}/${board}/${name}`
        pendingSet[k] = [Date.now() + setDelay, project, board, name, content]
        if (pendingInterval == null) {
            pendingInterval = setInterval(setPendingTasks, setDelay)
        }
    }

    static setTask(project, board, name, content) {
        name = encodeURIComponent(name)
        return axios.put(`/api/v1/projects/${project}/boards/${board}/${name}`, content, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
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
            .catch(loginWhenUnauthorized);
    }

    static deleteTask(project, board, name) {
        name = encodeURIComponent(name)
        return axios.delete(`/api/v1/projects/${project}/boards/${board}/${name}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }



    static touchTask(project, board, name) {
        name = encodeURIComponent(name)
        return axios.post(`/api/v1/projects/${project}/boards/${board}/${name}?touch`, null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized)
    }

    static uploadFileToLibrary(project, path, file = None) {
        path = encodeURIComponent(path)
        const config = getConfig();
        let formData = null;

        if (file) {
            formData = new FormData();
            formData.append("file", file);
            config.headers = config.headers || {};
            config.headers['Content-Type'] = 'multipart/form-data';
        }

        return axios.post(`/api/v1/projects/${project}/library${path}`, formData, config)
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static deleteFromLibrary(project, path) {
        path = encodeURIComponent(path)
        return axios.delete(`/api/v1/projects/${project}/library${path}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static downloadFromlibrary(project, path) {
        path = encodeURIComponent(path)
        //        return axios.get(`/api/v1/projects/${project}/library${path}`, getConfig())

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
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static createFolderInLibrary(project, path) {
        path = encodeURIComponent(path)
        return axios.put(`/api/v1/projects/${project}/library${path}`, null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static moveFileInLibrary(project, oldpath, path) {
        oldpath = encodeURIComponent(oldpath)
        path = encodeURIComponent(path)
        return axios.post(`/api/v1/projects/${project}/library${path}?move=${oldpath}`,
            null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static getSuggestions(project, prefix, total) {
        prefix = encodeURIComponent(prefix)
        let url = `/api/v1/projects/${project}/index/suggest/${prefix}`
        if (total) url += `&total=${total}`
        return axios.get(url, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static getGitStatus(project) {
        return axios.get(`/api/v1/projects/${project}/git/status`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

}

export default Server;