


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

export { lazyCall, sendPendingCalls };