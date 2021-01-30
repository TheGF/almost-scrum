
import { React, useContext, useState } from "react";
import ReactMde from "react-mde";
import "react-mde/lib/styles/css/react-mde-all.css";
import Server from '../server';
import UserContext from '../UserContext';
import MarkdownView from './MarkdownView';

function MarkdownEditor(props) {
    const { project } = useContext(UserContext);
    const { imageFolder, height, onSave, ...more } = props
    const [value, setValue] = useState(null)
    const [selectedTab, setSelectedTab] = useState("write");

    const saveImage = async function* (data) {
        const name = `${Date.now()}`
        await Server.uploadFileToLibrary(project, imageFolder, new Blob([data]), name)
        yield `/api/v1/projects/${project}/library${imageFolder}/${name}`;
        return true;
    };

    function onChange(value) {
        setValue(value)
        props.onChange(value);
    }

    function Image(props) {
        const token = localStorage.token
        if (token) {
            const src = `${props.src}?token=${token}`
            return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%' }} src={src} />
        } else {
            return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%' }} />
        }
    }


    console.log('Height', height)

    return <ReactMde
        key={height}
        value={value}
        onChange={onChange}
        minEditorHeight={height-120}
        maxEditorHeight={height-120}
        minPreviewHeight={height}
        selectedTab={selectedTab}
        onTabChange={setSelectedTab}
        generateMarkdownPreview={(markdown) =>
            Promise.resolve(<MarkdownView source={markdown} />)
        }
        paste={{
            saveImage: saveImage
        }}
        {...more}
    />
}
export default MarkdownEditor