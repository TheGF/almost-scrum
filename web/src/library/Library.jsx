import {
    Breadcrumb, BreadcrumbItem, BreadcrumbLink, Button, ButtonGroup,
    Editable, EditableInput, EditablePreview,
    HStack, Link, Spacer, Table, Tbody, Td, Th, Thead, Tr, VStack,
    IconButton, VisuallyHidden
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from 'react';
import { AiOutlineReload } from 'react-icons/ai';
import { BiEdit } from 'react-icons/bi';
import { CgRename } from 'react-icons/cg';
import { GoHome } from 'react-icons/go';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import Utils from '../core/utils';

function Library(props) {
    const { project } = useContext(UserContext);
    const [path, setPath] = useState('')

    const [favorites, setFavorites] = useState(
        (localStorage.getItem('ash-lib-favs') || '').split(',').filter(f => f)
    )
    const [files, setFiles] = useState([])

    function updateFavorites(path, oldPath) {
        if (path) {
            oldPath = oldPath || path
            const fa = [path, ...favorites.filter(f => f != oldPath)].slice(0, 5)
            localStorage.setItem('ash-lib-favs', fa)
            setFavorites(fa)
        }
    }

    function listFolder() {
        Server.listLibrary(project, path)
            .then(setFiles)
    }
    useEffect(listFolder, [path])

    function newFolder() {
        const cnt = files.filter(f => f.name.startsWith('new folder ')).length
        Server.createFolderInLibrary(project, `${path}/new folder ${cnt}`)
            .then(listFolder)
    }

    function deleteFile(file) {
        Server.deleteFromLibrary(project, `${path}/${file.name}`)
            .than(listFolder)
    }

    function uploadFile(evt) {
        const file = evt.target.files[0];
        file && Server.uploadFileToLibrary(project, path, file)
            .then(setFiles);
    }

    function onFileClick(file) {
        if (file.dir) {
            const p = [path, file.name].join('/')
            setPath(p)
            updateFavorites(p)
        } else {
            Server.downloadFromlibrary(project, `${path}/${file.name}`);
        }
    }

    function renameFile(file, name) {
        const p = `${path}/${name}`
        const o = `${path}/${file.name}`
        if (file.dir) {
            Server.moveFileInLibrary(project, o, p)
                .then(listFolder)
        } else {
            Server.moveFileInLibrary(project, o, p)
                .then(_ => updateFavorites(p, o))
                .then(listFolder)
        }
    }

    function renderFavs() {
        return favorites.map(p => {
            const label = p.split('/').pop()
            return <Button onClick={_ => setPath(p)} isActive={p == path}>
                {label}
            </Button>
        })
    }

    const rows = files.map(file => <Tr key={file.name}>
        <Td>
            <Editable defaultValue={file.name} isPreviewFocusable={false}
                onSubmit={name => renameFile(file, name)}>
                {({ isEditing, onEdit }) => (
                    <HStack spacing={2}>
                        <>
                            {Utils.fileIcon(file.dir, file.mime)}
                            <Link href="#" onClick={_ => {
                                if (!isEditing) onFileClick(file)
                            }}>
                                <EditablePreview
                                    style={{ cursor: 'pointer', color: 'blue' }} />
                            </Link>
                            <EditableInput maxWidth={'90%'} />
                        </>
                        <IconButton variant="outline" size="xs" icon={<BiEdit />} onClick={
                            onEdit
                        } />
                    </HStack>
                )}
            </Editable></Td>
        <Td>{Utils.getFriendlyDate(file.modTime)}</Td>
        <Td>{file.size}</Td>
        <Td>
            <ButtonGroup size="sm" spacing={2}>
                <Button>Links</Button>
                <Button onClick={_ => deleteFile(file)}>Delete</Button>
            </ButtonGroup>
        </Td>
    </Tr>)

    const folders = path.split('/')
    let breadcrumbs = folders.reduce((acc, folder) => {
        if (folder == '') return acc
        const pv = acc.length ? acc[acc.length - 1] : ''
        const p = `${pv}/${folder}`
        return [...acc, p]
    }, [])
    breadcrumbs = breadcrumbs.map((p, index) => <BreadcrumbItem>
        <BreadcrumbLink href="#" onClick={_ => setPath(p)}>
            {folders[1 + index]}
        </BreadcrumbLink>
    </BreadcrumbItem>)
    const breadcrumb = <Breadcrumb>
        <BreadcrumbItem>
            <BreadcrumbLink href="#" onClick={_ => setPath('')}><GoHome /></BreadcrumbLink>
        </BreadcrumbItem>
        {breadcrumbs}
    </Breadcrumb>

    let hiddenInput = null

    return <VStack w="90%" align="left" >
        <HStack w="90%" borderWidth="2" borderColor="gray">
            {breadcrumb}
            <Spacer />
            <ButtonGroup variant="outline"> {renderFavs()}
            </ButtonGroup>
            <Spacer />
            <VisuallyHidden>
                <input type="file"
                    ref={el => hiddenInput = el}
                    onChange={uploadFile} />
            </VisuallyHidden>
            <Button onClick={_ => hiddenInput.click()}>
                <T>Upload</T>
            </Button>
            <Button onClick={newFolder}>
                <T>New Folder</T>
            </Button>
            <IconButton onClick={listFolder}>
                <AiOutlineReload onclick={listFolder} />
            </IconButton>
        </HStack>
        <Table>
            <Thead>
                <Tr>
                    <Th>Name</Th>
                    <Th>Modified</Th>
                    <Th>Size</Th>
                    <Th>Actions</Th>
                </Tr>
            </Thead>
            <Tbody>
                {rows}
            </Tbody>
        </Table>
    </VStack>
}

export default Library