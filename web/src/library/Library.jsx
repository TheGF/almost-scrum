import {
    Breadcrumb, BreadcrumbItem, BreadcrumbLink, Flex, VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from 'react';
import { GoHome } from 'react-icons/go';
import Server from '../server';
import UserContext from '../UserContext';
import Bar from './Bar';
import Files from './Files';
import PageEditor from './PageEditor';

function Library(props) {
    const { project } = useContext(UserContext);
    const { attachedFiles, setAttachedFiles, height } = props
    const [path, setPath] = useState(props.path || '')

    const [favorites, setFavorites] = useState(
        (localStorage.getItem('ash-lib-favs') || '').split(',').filter(f => f)
    )
    const [files, setFiles] = useState([])
    const [page, setPage] = useState(null)

    function updateFavorites(path, oldPath) {
        if (path) {
            oldPath = oldPath || path
            const fa = [path, ...favorites.filter(f => f != oldPath)].slice(0, 4)
            localStorage.setItem('ash-lib-favs', fa)
            setFavorites(fa)
        }
    }

    function changePath(path) {
        if (path.endsWith('.pg')) {
            setPage(path)
        } else {
            setPath(path)
        }
        updateFavorites(path)
    }

    function getNextVersion(file) {
        const match = file.name.match(/(.*?)(~(\d+\.)+\d+)?(\.\w*)?$/)
        if (match.length < 5) return null

        const prefix = match[2] ? match[1] : `${match[1]}-`
        const ext = match[4] || ''

        const versions = files.map(file => {
            const match = file.name.match(/(.*?)(~(\d+\.)+\d+)?(\.\w*)?$/)
            if (match.length < 5) return null
            if (match[1] != prefix) return null
            return match[2]
        }).filter(v => v).sort()
        let version = null
        if (versions.length) {
            const last = versions[versions.length - 1]
            const match = last.match(/(~(\d+\.)+)(\d+)/)
            if (match.length == 4) {
                const last_digit = parseInt(match[3], 10)
                version = `${match[1]}${last_digit + 1}`
            }
        } 
        return [prefix, version, ext]
    }

    function addPageAttribute(f) {
        f.page = f.name.endsWith('.pg')
        f.dir = f.dir && !f.page
        return f
    }

    function listFolder() {
        Server.listLibrary(project, path)
            .then(files => setFiles(files.map(addPageAttribute)))
    }
    useEffect(listFolder, [path])

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

    return <VStack id="test-pdf" w="100%" align="left" >
        <PageEditor page={page} setPage={setPage} />
        <Bar path={path} changePath={changePath} getNextVersion={getNextVersion}
            favorites={favorites} files={files} listFolder={listFolder} />
        <Flex overflowY="auto"
            h={height && height - 160}>
            <Files path={path} changePath={changePath}
                files={files} listFolder={listFolder}
                getNextVersion={getNextVersion}
                attachedFiles={attachedFiles} setAttachedFiles={setAttachedFiles}
                updateFavorites={updateFavorites} />
        </Flex>
    </VStack>
}

export default Library