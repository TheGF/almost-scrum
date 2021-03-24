import React from 'react';
import {
    FaFile,
    FaFileArchive,
    FaFileAudio,
    FaFileExcel,
    FaFileImage,
    FaFilePdf,
    FaFilePowerpoint,
    FaFileVideo,
    FaFileWord,
    FaFolder
} from 'react-icons/fa';

const oneDay = 24 * 3600 * 1000;
let pendingCalls = {}


const mime2icons = [
    ['msword', FaFileWord]
    ['wordprocessingml', FaFileWord],
    ['ms-excel', FaFileExcel],
    ['spreadsheetml', FaFileExcel],
    ['powerpoint', FaFilePowerpoint],
    ['presentationml', FaFilePowerpoint],
    ['pdf', FaFilePdf],
    ['zip', FaFileArchive],
    ['vnd.rar', FaFileArchive],
    ['x-7z', FaFileArchive],
    ['image', FaFileImage],
    ['video', FaFileVideo],
    ['audio', FaFileAudio],
    ['video', FaFileImage],
    ['', FaFile]
];


class Utils {
    static getFriendlyDate(utcTime) {
        const unixtime = Date.parse(utcTime)
        const dateTime = new Date(unixtime);
        let today = new Date()
        today.setHours(0, 0, 0, 0)
        today = today.getTime();

        let date = null;
        let time = null;
        const locale = navigator.language || 'en-US';

        if (unixtime > today)
            date = 'today at';
        else if (unixtime > today - oneDay)
            date = 'yesterday at';
        else if (unixtime > today - 6 * oneDay)
            date = `${['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'][dateTime.getDay()]} at`;
        else {
            const options = { year: "2-digit", month: "2-digit", day: "2-digit" };
            date = dateTime.toLocaleDateString(locale, options);
        }

        time = dateTime.toLocaleTimeString(locale, { hour: '2-digit', minute: '2-digit' });
        return `${date} ${time}`;
    }

    static getFriendlySize(size) {
        if (size > 1024 * 1024 * 1024) {
            return `${Math.round(size/1024/1024/1024)}G`
        }
        if (size > 1024 * 1024) {
            return `${Math.round(size/1024/1024)}M`
        }
        if (size > 1024) {
            return `${Math.round(size/1024)}K`
        }
        return size
    }

    static lazyCall(key, action, callback = null) {
        pendingCalls[key] = {
            action: action,
            callback: callback
        };
    }

    static sendPendingCalls() {
        const calls = Object.values(pendingCalls);
        pendingCalls = {};
        for (const call of calls) {
            call.action();
            call.callback && call.callback();
        }
    }


    static fileIcon(dir, mime) {
        if (dir) return <FaFolder / > ;
        if (!mime) return <FaFile / > ;
        const Icon = mime2icons.filter(m => m && mime.includes(m[0]))[0][1]
        return <Icon / > ;
    }

    static isLocalhost() {
        return window.location.hostname === 'localhost' ||
            // [::1] is the IPv6 localhost address.
            window.location.hostname === '[::1]' ||
            // 127.0.0.1/8 is considered localhost for IPv4.
            window.location.hostname.match(
                /^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/
            )
    }

}

setInterval(Utils.sendPendingCalls, 30000);

export default Utils;