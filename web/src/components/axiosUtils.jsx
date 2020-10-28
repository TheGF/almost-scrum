
function getConfig() {
    const token = localStorage.token || '';

    return token ? {
        headers: {
            Authorization: `Bearer ${token}`,
        } 
    }: {}
}


function loginWhenUnauthorized(r) {
    if (r.response.status == 401) {
        localStorage.removeItem('username');
        localStorage.removeItem('token');
        window.location.assign(window.location.href);
    }
}

export { getConfig, loginWhenUnauthorized };