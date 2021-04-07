import { React, useEffect, useState } from "react";
import MarkdownEditor from '../core/MarkdownEditor';
import Server from '../server';


function HelpTab(props) {
    const { file } = props
    const [value, setValue] = useState('')

    function loadFile() {
        Server.getResource(file)
            .then(setValue)
    }
    useEffect(loadFile, [])

    return value && value.length && <MarkdownEditor value={value} readOnly={true} />

}
export default HelpTab