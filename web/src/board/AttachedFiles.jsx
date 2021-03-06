import {
    Button,
    ButtonGroup,
    Link,
    Table,
    Tbody,
    Td, Th, Thead, Tr, Flex
} from '@chakra-ui/react';
import { React, useContext, useState, useEffect } from "react";
import Utils from '../core/utils';
import Server from '../server';
import UserContext from '../UserContext';
import PageEditor from '../library/PageEditor';

function AttachedFiles(props) {
    const { attachedFiles, setAttachedFiles, readOnly, onShowPath, height } = props;
    const { project } = useContext(UserContext)
    const [infos, setInfos] = useState([])
    const [page, setPage] = useState(null)

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
        const p = `${info.parent}/${info.name}`.replace('//','/')
        const idx = attachedFiles.indexOf(p)
        const files = [
            ...attachedFiles.slice(0, idx),
            ...attachedFiles.slice(idx + 1)
        ]
        setAttachedFiles(files)
    }

    function onFileClick(info) {
        const p = `${info.parent}/${info.name}`
        if (info.name.endsWith(".pg")) {
            setPage(p)
        } else {
            Server.openFromlibrary(project, p);
        }
    }

    function onFolderClick(folder) {
        onShowPath(folder)
    }

    function deleteFile(info) {
        const p = `${info.parent}/${info.name}`
        Server.deleteFromLibrary(project, p)
            .then(_ => detach(info))
    }

    const rows = infos && infos.map(info => {
        const folder = '/' + info.parent.replace('.', '').replace(/^\/+|\/+$/g, '')

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
                    <Button onClick={_ => deleteFile(info)} isReadOnly={readOnly}>
                        Delete
                        </Button>
                </ButtonGroup>
            </Td>

        </Tr>
    })

    return <Flex overflowY="auto"
            h={height && height - 120}>
        <PageEditor page={page} setPage={setPage} />

        <Table>
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
    </Flex>
}
export default AttachedFiles;