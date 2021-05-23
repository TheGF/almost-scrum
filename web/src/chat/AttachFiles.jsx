
import React, { useState, Component } from 'react';
import { Badge, Button, HStack, IconButton, Input, Progress, Spacer, Tag, VisuallyHidden, VStack, Wrap, WrapItem } from "@chakra-ui/react";
import MicRecorder from 'mic-recorder-to-mp3';
import { BiTrash, BsFillPauseFill, BsFillPlayFill, BsFillStopFill, MdFiberManualRecord } from 'react-icons/all';
import Server from "../server";
import { render } from '@testing-library/react';
import T from '../core/T';


function AttachFiles(props) {
    const { value, onChange } = props
    const [files, setFiles] = useState(value.files)
    const [text, setText] = useState(value.text)
    let hiddenInput = null

    function addFiles(evt) {
        const fileList = evt.target.files
        const files = []
        for (let i = 0; i < fileList.length; i++) {
            files.push(fileList[i])
        }
        setFiles(files)
        onChange && onChange({text: text, files: files})
    }

    function changeText(evt) {
        const text = evt.target.value
        setText(text)
        onChange && onChange({text: text, files: files})
    }

    const filenames = files.map(file => <WrapItem>
        <Badge>{file.name}</Badge>
    </WrapItem>)
    return <VStack>
        <HStack width="100%">
            <Input size="sm" placeholder="Note or URL" value={text} onChange={changeText}></Input>
            <VisuallyHidden>
                <input type="file" multiple="multiple"
                    ref={el => hiddenInput = el}
                    onChange={addFiles} />
            </VisuallyHidden>
            <Button size="sm" onClick={_ => hiddenInput.click()} >
                <T>add</T>
            </Button>
        </HStack>
        <HStack width="100%"><Wrap width="100%">{filenames}</Wrap></HStack>
    </VStack>
}

export default AttachFiles