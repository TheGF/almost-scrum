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
import PageEditor from './PageEditor';
import { GiGrapes } from 'react-icons/gi';
import { MdCreateNewFolder } from 'react-icons/md';
import { VscNewFile } from 'react-icons/vsc';
import { Box } from '@chakra-ui/react';
import { Flex } from '@chakra-ui/react';

function Library(props) {
    const { project } = useContext(UserContext);
    const { attachedFiles, setAttachedFiles, height } = props
    const [path, setPath] = useState(props.path || '')
    const [uploading, setUploading] = useState(false)

    const [favorites, setFavorites] = useState(
        (localStorage.getItem('ash-lib-favs') || '').split(',').filter(f => f)
    )
    const [files, setFiles] = useState([])
    const [page, setPage] = useState(null)

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

    function newPage() {
        const cnt = files.filter(f => f.name.startsWith('page-')).length
        const folder = `${path}/page-${cnt}.pg`
        Server.createFolderInLibrary(project, folder)
            .then(listFolder)
            .then(_ => {
                const file = new Blob(['Change me'], { type: 'text/markdown' })
                Server.uploadFileToLibrary(project, folder, file, 'index.md')
            })
    }


    function deleteFile(file) {
        Server.deleteFromLibrary(project, `${path}/${file.name}`, file.name.endsWith('.pg'))
            .then(listFolder)
    }

    function uploadFile(evt) {
        setUploading(true)
        const file = evt.target.files[0];
        file && Server.uploadFileToLibrary(project, path, file)
            .then(setFiles)
            .then(_ => setUploading(false));
    }

    function onFileClick(file) {
        const p = [path, file.name].join('/')
        if (file.dir) {
            if (file.name.endsWith('.pg')) {
                setPage(p)
            } else {
                setPath(p)
                updateFavorites(p)
            }
        } else {
            Server.openFromlibrary(project, p);
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
            return <Button key={p} onClick={_ => setPath(p)} isActive={p == path}>
                {label}
            </Button>
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

    const rows = files.map(file => {
        const idx = file.name.lastIndexOf('.')
        const ext = idx != -1 ? file.name.substr(idx + 1) : ''
        const name = idx != -1 ? file.name.substr(0, idx) : file.name
        return <Tr key={file.name}>
            <Td>
                <Editable defaultValue={name} isPreviewFocusable={false}
                    onSubmit={name => renameFile(file, `${name}.${ext}`)}>
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
                                .{ext}
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
                    {getAttachButton(file)}
                    <Button onClick={_ => deleteFile(file)}>Delete</Button>
                </ButtonGroup>
            </Td>
        </Tr>
    })

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

    return <VStack w="100%" align="left" >
        <PageEditor page={page} setPage={setPage} />
        <HStack w="90%" borderWidth="2" borderColor="gray">
            {breadcrumb}
            <Spacer />
            <ButtonGroup variant="outline">
                {renderFavs()}
            </ButtonGroup>
            <Spacer />
            <VisuallyHidden>
                <input type="file"
                    ref={el => hiddenInput = el}
                    onChange={uploadFile} />
            </VisuallyHidden>
            <Button onClick={_ => hiddenInput.click()} isLoading={uploading} >
                <T>Upload</T>
            </Button>
            <Button onClick={newFolder} title="Create New Folder">
                <MdCreateNewFolder />
            </Button>
            <Button onClick={newPage} title="Create New Page">
                <VscNewFile />
            </Button>
            <IconButton onClick={listFolder} title="Reload">
                <AiOutlineReload onClick={listFolder} />
            </IconButton>
        </HStack>
        <Flex overflowY="auto"
            h={height && height - 160}>
            <Table w="100%">
                <Thead>
                    <Tr>
                        <Th>Name</Th>
                        <Th>Modified</Th>
                        <Th>Size</Th>
                        <Th>Actions</Th>
                    </Tr>
                </Thead>
                <Tbody  >
                    {rows}
                </Tbody>
            </Table>
        </Flex>
    </VStack>
}

export default Library