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

    static getProject(project) {
        return axios.get(`/api/v1/projects/${project}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static getUsers(project) {
        return axios.get(`/api/v1/projects/${project}/users`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }


    static getStoriesList(project, store) {
        return axios.get(`/api/v1/projects/${project}/stores/${store}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static getStory(project, store, story) {
        return axios.get(`/api/v1/projects/${project}/stores/${store}/${story}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static createStory(project, store, title, content) {
        return axios.post(`/api/v1/projects/${project}/stores/${store}?title=${title}`, content, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }


    static saveStory(project, store, story, content) {
        return axios.post(`/api/v1/projects/${project}/stores/${store}/${story}`, content, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static moveStory(project, store, story, target) {
        return axios.post(`/api/v1/projects/${project}/stores/${target}?from=${store}/${story}`,
            null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static touchStory(project, store, story) {
        return axios.post(`/api/v1/projects/${project}/stores/${store}/${story}?touch`, null, getConfig())
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
        link.href = `/api/v1/projects/${project}/library${path}`;
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