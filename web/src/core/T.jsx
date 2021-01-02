import React from 'react';

function translate(msg) {
    return msg ? (msg.charAt(0).toUpperCase() + msg.slice(1)) : '';
}

function T(props) {
    const msg = props.children
    return translate(msg);
}

T.translate = translate

export default T 

