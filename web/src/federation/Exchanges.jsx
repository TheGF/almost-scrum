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
    const [config, setConfig] = useState(null)
    const [status, setStatus] = useState(null)
    const toast = useToast()


    function getConfig() {
        Server.getFedConfig(project)
            .then(config => {
                const sked = {s3:[], webDAV:[], ftp:[], usb:[]}
                setConfig({...sked, ...config})
            })
        Server.getFedStatus(project)
            .then(setStatus)
    }
    useEffect(getConfig, [])

    function saveConfig() {
        Server.postFedConfig(project, config)
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

    let exchangesUI = config ? [
        ...config.s3.map(
            (exchange, i) => <S3Exchange exchange={exchange} update={v => updateExchange(config.s3, i, v)}
                status={status} />),
        ...config.webDAV.map(
            (exchange, i) => <WebDAVExchange exchange={exchange} update={v => updateExchange(config.webDAV, i, v)}
                connected={status && status.exchanges[exchange.name]} />),
        ...config.ftp.map(
            (exchange, i) => <FTPExchange exchange={exchange} update={v => updateExchange(config.ftp, i, v)}
                connected={status && status.exchanges[exchange.name]} />),
        ...config.usb.map(
            (exchange, i) => <USBExchange exchange={exchange} update={v => updateExchange(config.usb, i, v)}
                connected={status && status.exchanges[exchange.name]} />),
    ] : []

    function addS3() {
        config.s3.push({
            name: '',
            endpoint: '',
            bucket: '',
            accessKey: '',
            secret: '',
            useSSL: false,
            location: '',
        })
        setConfig({ ...config })
    }

    function addFTP() {
        config.ftp.push({
            name: '',
            url: '',
            username: '',
            password: '',
            secret: '',
            timeout: 10,
        })
        setConfig({ ...config })
    }


    function addUSBMedia() {
        config.usb.push({
        })
        setConfig({ ...config })
    }

    function addWebDAV() {
        config.webDAV.push({
            name: '',
            url: '',
            username: '',
            password: '',
            secret: '',
            timeout: 10,
        })
        setConfig({ ...config })
    }

    return config ? <VStack>
        <Table size="sm" padding="0" spacing="0">
            <Tr>
                <Td><FormLabel><T>Federation ID</T></FormLabel></Td>
                <Td><Input value={config.uuid} /></Td>
                <Td><FormLabel><T>span (days)</T></FormLabel></Td>
                <Td><Input type="number" value={config.span} /></Td>
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
            <Button colorScheme="blue" onClick={saveConfig}>Save</Button>
            <Button onClick={onClose}>Close</Button>
        </ButtonGroup>
        <Accordion w="100%" allowToggle>
            {exchangesUI}
        </Accordion>
    </VStack> : <span>Loading</span>
}
export default Exchanges;