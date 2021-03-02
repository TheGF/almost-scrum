import {
    Breadcrumb, BreadcrumbItem, BreadcrumbLink, Button, ButtonGroup,
    Editable, EditableInput, EditablePreview, HStack,
    IconButton, Link, Table, Tbody, Td, Th, Thead, Tr
} from '@chakra-ui/react';
import { React, useContext, useState } from 'react';
import { BiEdit } from 'react-icons/bi';
import { GiGrapes } from 'react-icons/gi';
import { GoHome } from 'react-icons/go';
import { GrUpgrade } from 'react-icons/gr';
import Utils from '../core/utils';
import Server from '../server';
import UserContext from '../UserContext';
import Versions from './Versions';

function Files(props) {
    const { project } = useContext(UserContext);
    const { path, changePath,
        files, listFolder,
        attachedFiles, setAttachedFiles,
        updateFavorites, getNextVersion } = props
    const isLocalhost = Utils.isLocalhost()
    const [showVersions, setShowVersions] = useState(null)

    function deleteFile(file) {
        Server.deleteFromLibrary(project, `${path}/${file.name}`, file.name.endsWith('.pg'))
            .then(listFolder)
    }

    function onFileClick(file) {
        const p = [path, file.name].join('/')
        if (file.dir || file.page) {
            changePath(p)
        } else {
            Server.openFromlibrary(project, p);
        }
    }

    function renameFile(file, name) {
        const p = `${path}/${name}`
        const o = `${path}/${file.name}`
        Server.moveFileInLibrary(project, o, p)
            .then(_ => {
                if (file.dir && !file.name.endsWith('.pg')) {
                    updateFavorites(p, o)
                }
                listFolder()
            })
    }

    function attach(file) {
        setAttachedFiles([...attachedFiles, file])
    }

    function detach(file) {
        const idx = attachedFiles.indexOf(file)
        const files = [
            ...attachedFiles.slice(0, idx),
            ...attachedFiles.slice(idx + 1)
        ]
        setAttachedFiles(files)
    }

    function getAttachButton(file) {
        if (!setAttachedFiles || file.dir && !file.name.endsWith('.pg')) {
            return null
        }
        const p = `${path}/${file.name}`
        if (attachedFiles.includes(p)) {
            return <Button colorScheme="yellow" onClick={_ => detach(p)}>Detach</Button>
        } else {
            return <Button onClick={_ => attach(p)}>Attach</Button>
        }
    }

    function openFile(file) {
        const p = `${path}/${file.name}`
        Server.localOpenFromLibrary(project, p)
    }

    function increaseVersion(file) {
        const [prefix, version, ext] = getNextVersion(file, true)
        const newName = version ?
            `${prefix}${version}${ext}` :
            `${prefix}-0.1${ext}`

        const p = `${path}/${newName}`
        const o = `${path}/${file.name}`
        Server.moveFileInLibrary(project, o, p)
            .then(listFolder)
    }

    function getOpenButton(file) {
        return isLocalhost ?
            <Button onClick={_ => openFile(file)}>Open</Button> :
            null
    }

    const rows = files && files.map(file => {
        const match = file.name.match(/(.*?)((\d+\.)+\d+)?(\.\w*)?$/)
        if (match.length < 5) {
            return ''
        }

        const name = match[1]
        const version = match[2] || ''
        const ext = match[4] || ''

        const versionUI = file.dir && !file.name.endsWith('.pg') ? null : version ?
            <HStack>
                <Link key="versions" onClick={_ => setShowVersions(file)}>
                    {version}
                </Link>
                <Link key="upgrade" onClick={_ => increaseVersion(file)}>
                    <GrUpgrade />
                </Link>
            </HStack> :
            <Link onClick={_ => increaseVersion(file)}>Enable</Link>

        return <Tr key={file.name}>
            <Td>
                <Editable defaultValue={name} isPreviewFocusable={false}
                    onSubmit={name => renameFile(file, `${name}${version}${ext}`)}>
                    {({ isEditing, onEdit }) => (
                        <HStack spacing={2}>
                            <>
                                {ext.endsWith('.pg') ? <GiGrapes /> : Utils.fileIcon(file.dir, file.mime)}
                                <Link href="#" onClick={_ => {
                                    if (!isEditing) onFileClick(file)
                                }}>
                                    <EditablePreview
                                        style={{ cursor: 'pointer', color: 'blue' }} />
                                </Link>
                                <EditableInput maxWidth={'90%'} />
                                {version}{ext}
                            </>
                            <IconButton variant="outline" size="xs" icon={<BiEdit />} onClick={
                                onEdit
                            } />
                        </HStack>
                    )}
                </Editable></Td>
            <Td>{Utils.getFriendlyDate(file.modTime)}</Td>
            <Td>{file.owner}</Td>
            <Td>{versionUI}</Td>
            <Td><span title={file.size}>{Utils.getFriendlySize(file.size)}</span></Td>
            <Td>
                <ButtonGroup size="sm" spacing={2}>
                    {getAttachButton(file)}
                    {getOpenButton(file)}
                    <Button onClick={_ => deleteFile(file)}>Delete</Button>
                </ButtonGroup>
            </Td>
        </Tr>
    })

    return <Table w="100%">
        <Versions path={path} file={showVersions} onClose={_ => setShowVersions(null)} />
        <Thead>
            <Tr>
                <Th>Name</Th>
                <Th>Modified</Th>
                <Th>Owner</Th>
                <Th>Version</Th>
                <Th>Size</Th>
                <Th>Actions</Th>
            </Tr>
        </Thead>
        <Tbody  >
            {rows}
        </Tbody>
    </Table>
}

export default Files