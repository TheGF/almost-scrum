import {
    Breadcrumb, BreadcrumbItem, BreadcrumbLink, Button, ButtonGroup,
    HStack, IconButton, Spacer, VisuallyHidden
} from '@chakra-ui/react';
import { React, useContext, useState } from 'react';
import { AiOutlineReload } from 'react-icons/ai';
import { GoHome } from 'react-icons/go';
import { MdCreateNewFolder } from 'react-icons/md';
import { VscNewFile } from 'react-icons/vsc';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import ConfirmUpload from './ConfirmUpload';

function Bar(props) {
    const { project } = useContext(UserContext);
    const { path, changePath,
        favorites, getNextVersion,
        files, listFolder } = props
    const [uploading, setUploading] = useState(false)
    const [confirmFileUpload, setConfirmFileUpload] = useState(null)

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

    function uploadFileToLibrary(file, name) {
        if (file) {
            Server.uploadFileToLibrary(project, path, file, name)
                .then(_ => {
                    listFolder()
                    setUploading(false)
                    setConfirmFileUpload(null)
                })
        } else {
            setUploading(false)
            setConfirmFileUpload(null)
        }
    }

    function uploadFile(evt) {
        const file = evt.target && evt.target.files[0];
        if (!file) return

        setUploading(true)
        const [_, version, __] = getNextVersion(file)

        if (version) {
            setConfirmFileUpload(file)
        } else {
            uploadFileToLibrary(file)
        }
    }

    function renderFavs() {
        return favorites.map(p => {
            const label = p.split('/').pop()
            return <Button key={p} onClick={_ => changePath(p)} isActive={p == path}>
                {label}
            </Button>
        })
    }

    function renderBreadcrumbs() {
        const folders = path.split('/')
        let breadcrumbs = folders.reduce((acc, folder) => {
            if (folder == '') return acc
            const pv = acc.length ? acc[acc.length - 1] : ''
            const p = `${pv}/${folder}`
            return [...acc, p]
        }, [])
        breadcrumbs = breadcrumbs.map((p, index) => <BreadcrumbItem>
            <BreadcrumbLink href="#" onClick={_ => changePath(p)}>
                {folders[1 + index]}
            </BreadcrumbLink>
        </BreadcrumbItem>)
        return <Breadcrumb>
            <BreadcrumbItem>
                <BreadcrumbLink href="#" onClick={_ => changePath('')}>
                    <GoHome />
                </BreadcrumbLink>
            </BreadcrumbItem>
            {breadcrumbs}
        </Breadcrumb>

    }

    let hiddenInput = null

    return <HStack w="90%" borderWidth="2" borderColor="gray">
        <ConfirmUpload file={confirmFileUpload}
            uploadFileToLibrary={uploadFileToLibrary}
            getNextVersion={getNextVersion} />
        {renderBreadcrumbs()}
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
}

export default Bar