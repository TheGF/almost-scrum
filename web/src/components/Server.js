import axios from 'axios';
import { getConfig, loginWhenUnauthorized } from './axiosUtils';

class Server {

    static getProjectsList() {
        return axios.get('/api/v1/projects', getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static getStoriesList(project, store) {
        return axios.get(`/api/v1/projects/${project}/${store}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static getStory(project, store, story) {
        return axios.get(`/api/v1/projects/${project}/${store}/${story}`, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }


    static createStory(project, store, title) {
        return axios.post(`/api/v1/projects/${project}/${store}?title=${title}`, emptyStory, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static saveStory(project, store, story, content) {
        return axios.post(url, content, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);        
    }

    static moveStory(project, store, story, target) {
        return axios.post(`/api/v1/projects/${project}/${target}?from=${store}/${story}`,
            null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized);
    }

    static touchStory(project, store, story) {
        return axios.post(`/api/v1/projects/${project}/${store}/${story}?touch`, null, getConfig())
            .then(r => r.data)
            .catch(loginWhenUnauthorized)
    }





}

export default Server;