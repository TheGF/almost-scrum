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

class Server {

    static getProjectsList() {
        return axios.get('/api/v1/projects', getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    // static getProject(project) {
    //     return axios.get(`/api/v1/projects/${project}`, getConfig())
    //         .then(r => r.data)
    //         .catch(loginWhenUnauthorized);
    // }

    static getUsers(project) {
        return axios.get(`/api/v1/projects/${project}/users`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static listBoards(project) {
        return axios.get(`/api/v1/projects/${project}/boards`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static listTasks(project, board, filter, start, end) {
        let url = `/api/v1/projects/${project}/boards/${board}?`
        if (start) url += `start=${start}`;
        if (end) url += `end=${end}`;
        if (filter) url += `filter=${filter}`;

        return axios.get(url, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static getTask(project, board, name) {
        return axios.get(`/api/v1/projects/${project}/boards/${board}/${name}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }


    static createTask(project, board, title, content) {
        return axios.post(`/api/v1/projects/${project}/boards/${board}?title=${title}`, content, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }


    static setTask(project, board, name, content) {
        return axios.put(`/api/v1/projects/${project}/boards/${board}/${name}`, content, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static moveTask(project, board, source) {
        return axios.post(`/api/v1/projects/${project}/boards/${board}?from=${source}`,
            null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static touchTask(project, board, name) {
        return axios.post(`/api/v1/projects/${project}/boards/${board}/${name}?touch`, null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized)
    }

    static libraryPost(project, path, file = None) {
        const config = getConfig();
        let formData = null;

        if (file) {
            formData = new FormData();
            formData.append("file", file);
            config.headers['Content-Type'] = 'multipart/form-data';
        }

        return axios.post(`/api/v1/projects/${project}/library${path}`, formData, config)
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static libraryDelete(project, path) {
        return axios.delete(`/api/v1/projects/${project}/library${path}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static libraryDownload(project, path) {
//        return axios.get(`/api/v1/projects/${project}/library${path}`, getConfig())

        const link = document.createElement("a");
        link.href = `/api/v1/projects/${project}/library${path}?token=${localStorage.token}`;
        link.target = '_';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    static libraryList(project, path) {
        return axios.get(`/api/v1/projects/${project}/library${path}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

}

export default Server;