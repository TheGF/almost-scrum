import {
    Accordion, Button, ButtonGroup, FormLabel, Input, Spacer, Table, Td, Tr, useToast, VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import FTPExchange from './exchanges/FTPExchange';
import S3Exchange from './exchanges/S3Exchange';
import WebDAVExchange from './exchanges/WebDAVExchange';
import USBExchange from './exchanges/USBExchange';


function Exchanges(props) {
    const { project } = useContext(UserContext)
    const { onClose } = props
    const [transport, setTransport] = useState(null)
    const [status, setStatus] = useState(null)
    const toast = useToast()


    function getConfig() {
        Server.getFedTransport(project)
            .then(setTransport)
        Server.getFedState(project)
            .then(setStatus)
    }
    useEffect(getConfig, [])

    function saveTransport() {
        Server.postFedTransport(project, transport)
            .then(_ => toast({
                title: `Config Saved`,
                description: 'The federation configuration has been saved',
                status: "success",
                duration: 9000,
                isClosable: true,
            }))
    }

    function updateExchange(l, idx, value) {
        if (value) {
            l[idx] = value
        } else {
            l.splice(idx, 1)
            setConfig({ ...config })
        }
    }

    let exchangesUI = transport ? [
        ...transport.s3.map(
            (exchange, i) => <S3Exchange exchange={exchange} update={v => updateExchange(transport.s3, i, v)}
                status={status} />),
        ...transport.webDAV.map(
            (exchange, i) => <WebDAVExchange exchange={exchange} update={v => updateExchange(transport.webDAV, i, v)}
                connected={status && status.netStats[exchange.name]} />),
        ...transport.ftp.map(
            (exchange, i) => <FTPExchange exchange={exchange} update={v => updateExchange(transport.ftp, i, v)}
                connected={status && status.netStats[exchange.name]} />),
        ...transport.usb.map(
            (exchange, i) => <USBExchange exchange={exchange} update={v => updateExchange(transport.usb, i, v)}
                connected={status && status.netStats[exchange.name]} />),
    ] : []

    function addS3() {
        transport.s3.push({
            name: '',
            endpoint: '',
            bucket: '',
            accessKey: '',
            secret: '',
            useSSL: false,
            location: '',
        })
        setTransport({ ...transport })
    }

    function addFTP() {
        transport.ftp.push({
            name: '',
            url: '',
            username: '',
            password: '',
            secret: '',
            timeout: 10,
        })
        setTransport({ ...transport })
    }


    function addUSBMedia() {
        transport.usb.push({
        })
        setTransport({ ...transport })
    }

    function addWebDAV() {
        transport.webDAV.push({
            name: '',
            url: '',
            username: '',
            password: '',
            secret: '',
            timeout: 10,
        })
        setTransport({ ...transport })
    }

    return transport ? <VStack>
        <Table size="sm" padding="0" spacing="0">
            <Tr>
                <Td><FormLabel><T>Federation ID</T></FormLabel></Td>
                <Td><Input value={transport.fedId} /></Td>
                <Td><FormLabel><T>span (days)</T></FormLabel></Td>
                <Td><Input type="number" value={transport.span} /></Td>
            </Tr>
            <Tr>
            </Tr>
        </Table>
        <Spacer minHeight="1em" />
        <ButtonGroup>
            <Button onClick={_ => addS3()}><T>add S3</T></Button>
            <Button onClick={_ => addWebDAV()}><T>add WebDAV</T></Button>
            <Button onClick={_ => addFTP()}><T>add FTP</T></Button>
            <Button onClick={_ => addUSBMedia()}><T>add USB Media</T></Button>
            <Spacer minW="2em" />
            <Button colorScheme="blue" onClick={saveTransport}>Save</Button>
            <Button onClick={onClose}>Close</Button>
        </ButtonGroup>
        <Accordion w="100%" allowToggle>
            {exchangesUI}
        </Accordion>
    </VStack> : <span>Loading</span>
}
export default Exchanges;