import {
    Breadcrumb, BreadcrumbItem, BreadcrumbLink, Button, ButtonGroup,
    HStack, IconButton, Spacer, VisuallyHidden
} from '@chakra-ui/react';
import { React, useContext, useState } from 'react';
import { AiOutlineReload } from 'react-icons/ai';
import { GiBlackBook } from 'react-icons/gi';
import { GoHome } from 'react-icons/go';
import { MdCreateNewFolder } from 'react-icons/md';
import { VscNewFile } from 'react-icons/vsc';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import ConfirmUpload from './ConfirmUpload';
import Book from './Book';
//import html2pdf from 'html2pdf.js';

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
        const name = `page-${cnt}.pg`
        const file = new Blob(['Change me'], { type: 'text/markdown' })
        Server.uploadFileToLibrary(project, path, file, name)
            .then(listFolder)

    }

    function newBook() {
        // doc.fromHTML(document.getElementById("test-pdf"), // page element which you want to print as PDF
        //     15,
        //     15,
        //     {
        //         'width': 170  //set width
        //     },
        //     function (a) {
        //         doc.save("HTML2PDF.pdf"); // save file name as HTML2PDF.pdf
        //     })

        Server.postNewBook2(project, path)
            .then(html => {
                const opt = {
                    margin: [10, 10, 10, 10],
                    filename: `document.pdf`,
                    image: {
                        type: 'jpeg',
                        quality: 0.98
                    },
                    html2canvas: {
                        scale: 2,
                        useCORS: false
                    },
                    jsPDF: {
                        unit: 'mm',
                        format: 'letter',
                        orientation: 'portrait'
                    }
                };
                html2pdf().from(html).set(opt).save();

                // const pdf = new jsPDF('p', 'pt', 'letter');;  //create jsPDF object
                // const margin = [40, 80, 40, 80]

                // const body = html.match(/<body[^>]*>([\w|\W]*)<\/body>/im)[0];

                // pdf.html(html, // HTML string or DOM elem ref.
                //     {
                //         margin: margin,
                //         callback: _ => pdf.save('Test.pdf')
                //     }
                // );

                //body = document.getElementById('document');
                // const body = html.match(/<body[^>]*>([\w|\W]*)<\/body>/im)[0];

                // const anchor = document.createElement("div");
                // anchor.innerHTML = body

                // var opt = {
                //     margin:       [10, 0, 10, 0],
                //     filename:     `document.pdf`,
                //     image:        { type: 'jpeg', quality: 0.98 },
                //     html2canvas:  { scale: 2, useCORS: false },
                //     jsPDF:        { unit: 'mm', format: 'A4', orientation: 'portrait' }
                // };
                // html2pdf().from(anchor).set(opt).save();
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
            <BreadcrumbLink key={index} href="#" onClick={_ => changePath(p)}>
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
        <Book path={path} listFolder={listFolder}/>
        <IconButton onClick={listFolder} title="Reload">
            <AiOutlineReload onClick={listFolder} />
        </IconButton>
    </HStack>
}

export default Bar