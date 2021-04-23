
import '@toast-ui/editor/dist/toastui-editor.css';
import { Editor, Viewer } from '@toast-ui/react-editor';
import 'codemirror/lib/codemirror.css';
import { React, useContext, useRef, useState } from "react";
import Server from '../server';
import UserContext from '../UserContext';
import MarkdownImage from './MarkdownImage';
import uml from '@toast-ui/editor-plugin-uml';
import tableMergedCell from '@toast-ui/editor-plugin-table-merged-cell';
import 'tui-color-picker/dist/tui-color-picker.css';
import colorSyntax from '@toast-ui/editor-plugin-color-syntax';
import 'highlight.js/styles/github.css';
import codeSyntaxHighlight from '@toast-ui/editor-plugin-code-syntax-highlight';
import 'tui-chart/dist/tui-chart.css';
import chart from '@toast-ui/editor-plugin-chart';

function MarkdownEditor(props) {
    const { project } = useContext(UserContext);
    const { imageFolder, height, onSave, disablePreview, readOnly, ...more } = props
    const projectPath = `/api/v1/projects/${project}`

    const [value, setValue] = useState(props.value && `${props.value}`
                                        .replaceAll('~', projectPath) || null);
    const editorRef = useRef(null)
    const [editImage, setEditImage] = useState(null)
    const [refresh, setRefresh] = useState(false)
    
    function uploadFile(e) {
        function replaceHtml() {
            let url = `/api/v1/projects/${project}/library`
            url += `${imageFolder}/${name}#size=50,align=center`
            target.innerHTML = MarkdownImage.getImg(url)
        }

        const name = `${Date.now()}`
        const items = e.clipboardData && e.clipboardData.items
        const target = e.target

        if (items) {
            for (const item of items) {
                if (item.kind == 'file') {
                    const blob = item.getAsFile()
                    Server.uploadFileToLibrary(project, imageFolder,
                        blob, name)
                        .then(_=> {
                            replaceHtml()
                            Server.setVisibility(project, `${imageFolder}/${name}`, true)
                        })
                }
            }
        }
    }

    function onLoad(editor) {
        editor.getUI().el.addEventListener('paste', uploadFile)
    }

    function onChange() {
        const v = editorRef.current.getInstance().getMarkdown()
        if (!v) return
        
        if (!value || v.trim() != value.trim()) {
            console.log('SAVING')
            setValue(v)
            props.onChange(v.replaceAll(projectPath, '~'));
        }
    }

    function addImageBlobHook(blob, callback) {
        const name = `${Date.now()}`
        let url = `/api/v1/projects/${project}/library`
        url += `${imageFolder}/${name}#size=50,align=center`

        Server.uploadFileToLibrary(project, imageFolder,
            blob, name)
            .then(_=> {
                callback(url, ' ')
                Server.setVisibility(project, `${imageFolder}/${name}`, true)
                setTimeout(_=>setRefresh(!refresh), 500)
            })
        return false;
    }

    function renderImage(node, context) {
        function setOpenImageSettings() {
            const img = document.getElementById(id)
            if (img) {
                img.onclick = _ => setEditImage(img)
            } else {
                setTimeout(setOpenImageSettings, 20)
            }
        }

        if (!context.entering) return null

        const id = `${Date.now()}.${Math.random()}`
        const content = MarkdownImage.getImg(node.destination, id)
        setTimeout(setOpenImageSettings, 20)

        return {
            type: 'html',
            content: content,
        };
    }

    const plugins = [chart, codeSyntaxHighlight, colorSyntax, tableMergedCell, uml]

    const gap = disablePreview ? 120 : 20
    return readOnly ? <Viewer
        initialValue={value}
        previewStyle="tab"
        height={height - 40}
        initialEditType="wysiwyg"
        useCommandShortcut={true}
        ref={editorRef}
        customHTMLRenderer={{
            image: renderImage,
        }}
        plugins={plugins}
    /> : <>
        <MarkdownImage image={editImage} setImage={setEditImage} />
        <Editor
            key={refresh}
            initialValue={value}
            previewStyle="tab"
            height={height - 40}
            initialEditType="wysiwyg"
            useCommandShortcut={true}
            ref={editorRef}
            hooks={{
                addImageBlobHook: addImageBlobHook,
            }}
            events={{
                load: onLoad,
                change: onChange,
            }}
            customHTMLRenderer={{
                image: renderImage,
            }}
            plugins={plugins}
        />
    </>

}
export default MarkdownEditor