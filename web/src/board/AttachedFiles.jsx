import {
    Button,
    ButtonGroup,
    Link,
    Table,
    Tbody,
    Td, Th, Thead, Tr
} from '@chakra-ui/react';
import { React, useContext, useState, useEffect } from "react";
import Utils from '../core/utils';
import Server from '../server';
import UserContext from '../UserContext';

function AttachedFiles(props) {
    const { attachedFiles, setAttachedFiles, readOnly, onShowPath } = props;
    const { project, info } = useContext(UserContext)
    const [infos, setInfos] = useState([])

    function getStat() {
        if (attachedFiles) {
            Server.getLibraryStat(project, attachedFiles)
                .then(setInfos)
        } else {
            setInfos([])
        }
    }
    useEffect(getStat, [attachedFiles])

    function detach(info) {
        const p = `${info.parent}/${info.name}`
        const idx = attachedFiles.indexOf(p)
        setAttachedFiles([
            ...attachedFiles.slice(0, idx),
            ...attachedFiles.slice(idx + 1)
        ])
    }

    function onFileClick(info) {
        Server.downloadFromlibrary(project, `${info.parent}/${info.name}`);
    }

    function onFolderClick(folder) {
        onShowPath(folder)
    }


    const rows = infos && infos.map(info => {
        const folder = '/'+info.parent.replace('.', '').replace(/^\/+|\/+$/g, '')

        return <Tr key={`${info.parent}/${info.name}`}>
            <Td><Link onClick={_ => onFileClick(info)}>{info.name}</Link></Td>
            <Td><Link onClick={_ => onFolderClick(folder)}>{folder}</Link></Td>
            <Td>{Utils.getFriendlyDate(info.modTime)}</Td>
            <Td>{info.size}</Td>
            <Td>
                <ButtonGroup size="sm" spacing={2}>
                    <Button onClick={_ => detach(info)} isReadOnly={readOnly}>
                        Detach
                        </Button>
                    <Button onClick={_ => deleteFile(file)} isReadOnly={readOnly}>
                        Delete
                        </Button>
                </ButtonGroup>
            </Td>

        </Tr>
    })

    return <Table>
        <Thead>
            <Tr>
                <Th>Name</Th>
                <Th>Folder</Th>
                <Th>Modified</Th>
                <Th>Size</Th>
                <Th>Actions</Th>
            </Tr>
        </Thead>
        <Tbody>
            {rows}
        </Tbody>
    </Table>
}
export default AttachedFiles;