
function getConfig() {
    const token = localStorage.token || '';

    return token ? {
        headers: {
            Authorization: `Bearer ${token}`,
        } 
    }: {}
}


function loginWhenUnauthorized(r) {
    if (r && r.response && r.response.status == 401) {
        localStorage.removeItem('username');
        localStorage.removeItem('token');
        window.location.assign(window.location.href);
    }
}

let pendingCalls = {}

function lazyCall(key, action, callback = null) {
    pendingCalls[key] = {
        action: action,
        callback: callback
    };
}

function sendPendingCalls() {
    const calls = Object.values(pendingCalls);
    pendingCalls = {};
    for (const call of calls) {
        call.action();
        call.callback && call.callback();
    } 
}

setInterval(sendPendingCalls, 30000);

export { getConfig, loginWhenUnauthorized, lazyCall, sendPendingCalls };