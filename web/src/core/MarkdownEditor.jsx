
import { React, useContext, useState } from "react";
import ReactMde from "react-mde";
import "react-mde/lib/styles/css/react-mde-all.css";
import Server from '../server';
import UserContext from '../UserContext';
import MarkdownView from './MarkdownView';

function MarkdownEditor(props) {
    const { project } = useContext(UserContext);
    const { imageFolder, height, onSave, disablePreview, ...more } = props
    const [value, setValue] = useState(null)
    const [selectedTab, setSelectedTab] = useState(disablePreview ? "write" : "preview");

    const saveImage = async function* (data) {
        const name = `${Date.now()}`
        await Server.uploadFileToLibrary(project, imageFolder, new Blob([data]), name)
        yield `~library${imageFolder}/${name}`;
        return true;
    };

    function onChange(value) {
        setValue(value)
        props.onChange(value);
    }
    

    return <ReactMde
        key={height}
        value={value}
        onChange={onChange}
        disablePreview={disablePreview}
        minEditorHeight={height - 120}
        maxEditorHeight={height - 120}
        minPreviewHeight={height}
        selectedTab={selectedTab}
        onTabChange={setSelectedTab}
        generateMarkdownPreview={(markdown) =>
            Promise.resolve(<div onClick={_=>setSelectedTab("write")}>
                <MarkdownView value={markdown}
                onChange={onChange} />
                </div>)
        }
        paste={{
            saveImage: saveImage
        }}
        {...more}
    />
}
export default MarkdownEditor