
import {
    Button, ButtonGroup, Input, Modal, ModalBody, ModalCloseButton,
    ModalContent, ModalFooter, ModalHeader, ModalOverlay, Radio, RadioGroup,
    Select, Stack, Switch, Table, Tbody, Td, Tr, useDisclosure, useToast, VStack
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import { GiBlackBook } from 'react-icons/gi';
import T from "../core/T";
import Server from '../server';
import UserContext from '../UserContext';
import html2pdf from './html2pdf.bundle.min'

function binToBlob(byteCharacters, contentType = '', sliceSize = 512) {
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += sliceSize) {
        const slice = byteCharacters.slice(offset, offset + sliceSize);

        const byteNumbers = new Array(slice.length);
        for (let i = 0; i < slice.length; i++) {
            byteNumbers[i] = slice.charCodeAt(i);
        }

        const byteArray = new Uint8Array(byteNumbers);
        byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: contentType });
    return blob;
}



function Book(props) {
    const { project } = useContext(UserContext);
    const { path, listFolder } = props
    const { isOpen, onOpen, onClose } = useDisclosure()
    const [styles, setStyles] = useState([
        'booklet', 'portrait', 'A4', 'pdf'
    ])
    const [title, setTitle] = useState(path.substring(path.lastIndexOf('/') + 1) || 'Library')
    const [subtitle, setSubtitle] = useState('')
    const [authors, setAuthors] = useState('')
    const [createInProgress, setCreateInProgress] = useState(false)
    const toast = useToast()

    function loadBookConfig() {
        Server.downloadFromlibrary(project, `${path}/.book.json`)
            .then(c => {
                if (c) {
                    if (c.title) setTitle(c.title)
                    if (c.styles) setStyles(c.styles)
                    setAuthors(c.authors || '')
                }
            })
    }
    //    useEffect(loadBookConfig, [])

    function saveBookConfig() {
        const c = {
            title: title,
            subtitle: subtitle,
            authors: authors,
            styles: styles,
        }
        Server.uploadFileToLibrary(project, path, new Blob([JSON.stringify(c)]), '.book.json')
    }

    function getStyle(options) {
        return styles.filter(s => options.includes(s))[0]
    }

    function flipStyle(style) {
        if (styles.contains(style)) {
            setStyles(styles.filter(s => s != style))
        } else {
            setStyles([...styles, style])
        }
    }

    function setStyle(options, style) {
        setStyles([...styles.filter(s => !options.includes(s)), style])
    }

    function convertToPdf(html) {
        const format = getStyle(papers) || 'A4'
        const orientation = getStyle(orientations) || 'landscape'

        const opt = {
            margin: [10, 10, 10, 10],
            filename: `${title}.pdf`,
            image: {
                type: 'jpeg',
                quality: 0.98
            },
            html2canvas: {
                scale: 2,
                useCORS: true,
            },
            jsPDF: {
                unit: 'mm',
                format: format,
                orientation: orientation,
            }
        };
        html2pdf().from(html).set(opt).toPdf().output()
            .then(pdf => {
                const blob = binToBlob(pdf, 'application/octet-stream')
                Server.uploadFileToLibrary(project, path, blob, `${title}.pdf`)
                    .then(_ => confirm(`${title}.pdf`))
            });
    }

    function confirm(name) {
        onClose()
        setCreateInProgress(false)
        toast({
            title: "Created!",
            description: `Saved in document ${name} in the folder`,
            status: "success",
            isClosable: true,
        })
        listFolder()
    }

    function createBook() {
        setCreateInProgress(true)
        const settings = {
            title: title,
            subtitle: subtitle,
            authors: authors,
            styles: styles,
        }
        Server.postNewBook2(project, path, settings)
            .then(html => {
                const format = getStyle(formats)
                if (format == 'pdf') {
                    convertToPdf(html)
                }
                if (format == 'html') {
                    Server.uploadFileToLibrary(project, path, new Blob([html]),
                        `${title}.html`)
                        .then(_ => confirm(`${title}.html`))
                }
                saveBookConfig()
            })
    }

    const kinds = ['booklet', 'article', 'paper']
    const kindUI = <RadioGroup value={getStyle(kinds)} colorScheme="blue"
        onChange={s => setStyle(kinds, s)} >
        <Stack spacing={4} direction="row">
            <Stack spacing={4} direction="row">{
                kinds.map(v => <Radio value={v}><T>{v}</T></Radio>)
            }</Stack>
        </Stack>
    </RadioGroup>

    const orientations = ['portrait', 'landscape']
    const orientationUI = <RadioGroup value={getStyle(orientations)} colorScheme="blue"
        onChange={s => setStyle(orientations, s)} >
        <Stack spacing={4} direction="row">{
            orientations.map(v => <Radio value={v}><T>{v}</T></Radio>)
        }</Stack>
    </RadioGroup>

    const papers = ['A4', 'letter', 'A5', 'A3', 'B5', 'B4']
    const paperUI = <Select placeholder="Select option" value={getStyle(papers)}
        size="sm" onChange={e => setStyle(papers, e.target.value)} >{
            papers.map(v => <option value={v}>{v}</option>)
        }</Select>

    const titleUI = <Input size="sm" value={title}
        onChange={e => setTitle(e.target.value)} />
    const subTitleUI = <Input size="sm" value={subtitle}
        onChange={e => setSubtitle(e.target.value)} />
    const authorsUI = <Input size="sm" value={authors}
        onChange={e => setAuthors(e.target.value)} />

    const tocUI = <Switch value={getStyle(['toc'])} onChange={_ => flipStyle('toc')} />
    const numberedHeadersUI = <Switch value={getStyle(['numberedHeaders'])} onChange={_ => flipStyle('numberedHeaders')} />

    const formats = ['pdf', 'html']
    const formatUI = <RadioGroup value={getStyle(formats)} colorScheme="blue"
        onChange={s => setStyle(formats, s)} >
        <Stack spacing={4} direction="row">{
            formats.map(v => <Radio value={v}><T>{v}</T></Radio>)
        }</Stack>
    </RadioGroup>

    return <>
        <Button onClick={onOpen} title="Create a Book">
            <GiBlackBook />
        </Button>

        <Modal isOpen={isOpen} onClose={onClose} size="lg">
            <ModalOverlay />
            <ModalContent >
                <ModalHeader>Create a Document</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <VStack>
                        <Table spacing={1} size="sm">
                            <Tbody>
                                <Tr>
                                    <Td><T>kind</T></Td>
                                    <Td>{kindUI}</Td>
                                </Tr>
                                <Tr>
                                    <Td><T>orientation</T></Td>
                                    <Td>{orientationUI}</Td>
                                </Tr>
                                <Tr>
                                    <Td><T>paper</T></Td>
                                    <Td>{paperUI}</Td>
                                </Tr>
                                <Tr>
                                    <Td><T>title</T></Td>
                                    <Td>{titleUI}</Td>
                                </Tr>
                                <Tr>
                                    <Td><T>subtitle</T></Td>
                                    <Td>{subTitleUI}</Td>
                                </Tr>
                                <Tr>
                                    <Td><T>author(s)</T></Td>
                                    <Td>{authorsUI}</Td>
                                </Tr>
                                <Tr>
                                    <Td><T>table of contents</T></Td>
                                    <Td>{tocUI}</Td>
                                </Tr>
                                <Tr>
                                    <Td><T>numbered headers</T></Td>
                                    <Td>{numberedHeadersUI}</Td>
                                </Tr>
                                <Tr>
                                    <Td><T>format</T></Td>
                                    <Td>{formatUI}</Td>
                                </Tr>
                            </Tbody>
                        </Table>
                        <ButtonGroup>
                            <Button onClick={createBook} colorScheme="blue"
                                isLoading={createInProgress}>
                                Create
                            </Button>
                            <Button onClick={onClose}>Close</Button>
                        </ButtonGroup>
                    </VStack>
                </ModalBody>
                <ModalFooter>
                </ModalFooter>
            </ModalContent>
        </Modal>
    </>
}
export default Book