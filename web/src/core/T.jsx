import React from 'react';

function T(props) {
    const msg = props.children
    return msg ? (msg.charAt(0).toUpperCase() + msg.slice(1)) : '';
}
export default T

