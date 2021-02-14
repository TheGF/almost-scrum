import { React, useEffect, useState } from "react";
import MarkdownView from '../core/MarkdownView';
import Server from '../server';


function HelpTab(props) {
    const { file } = props
    const [value, setValue] = useState('')

    function loadFile() {
        Server.getResource(file)
            .then(setValue)
    }
    useEffect(loadFile, [])

    return <div className="mde-preview">
        <div className="mde-preview-content task-viewer" style={{ overflow: 'auto' }}>
            <MarkdownView value={value} readOnly={true} />
        </div>
    </div>
}
export default HelpTab