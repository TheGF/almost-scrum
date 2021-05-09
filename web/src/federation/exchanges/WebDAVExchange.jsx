import {
    AccordionButton,
    AccordionIcon,
    AccordionItem,
    AccordionPanel, Box, Button, Input, Switch, Table, Tag, Td, Tr
} from '@chakra-ui/react';
import { React, useState } from "react";
import T from '../../core/T';



function WebDAVExchange(props) {
    const [exchange, setExchange] = useState(props.exchange)
    const connectErr = status && status.connectErr    
    const upload = status && Utils.getFriendlySize(status.throughputUp)
    const download = status && Utils.getFriendlySize(status.throughputDown)

    function updateField(evt, field, value) {
        exchange[field] = value != undefined ? value : evt.target.value
        const u = { ...exchange }
        setExchange(u)
        props.update(u)
    }

    return <AccordionItem>
        <h2>
            <AccordionButton>
                <Box flex="1" textAlign="left">
                    WebDAV - {exchange.name}
                </Box>
                {connectErr == null ?
                    <Tag colorScheme="green" title={`U:${upload}-D:${download}`}><T>connected</T></Tag> :
                    <Tag colorScheme="red">{connectErr && connectErr.Msg ? connectErr.Msg : <T>disconnected</T> }</Tag>}
                <AccordionIcon />
            </AccordionButton>
        </h2>
        <AccordionPanel w="100%">
            <Table size="sm" padding="0">
                <Tr>
                    <Td><label><T>name</T></label></Td>
                    <Td><Input type="text" value={exchange.name} onChange={e => updateField(e, 'name')} /></Td>
                    <Td><Button w="100%" onClick={_=>props.update(null)}>Delete</Button></Td>
                </Tr>
                <Tr>
                    <Td><label><T>url</T></label></Td>
                    <Td><Input type="text" value={exchange.url} onChange={e => updateField(e, 'url')} /></Td>
                </Tr>
                <Tr>
                    <Td><label><T>username</T></label></Td>
                    <Td><Input type="text" value={exchange.username} onChange={e => updateField(e, 'username')} /></Td>
                </Tr>
                <Tr>
                    <Td><label><T>password</T></label></Td>
                    <Td><Input type="text" value={exchange.password} onChange={e => updateField(e, 'password')} /></Td>
                </Tr>
                <Tr>
                    <Td><label><T>timeout</T></label></Td>
                    <Td><Input type="password" value={exchange.timeout} onChange={e => updateField(e, 'timeout')} /></Td>
                </Tr>
            </Table>

        </AccordionPanel>
    </AccordionItem>
}
export default WebDAVExchange;