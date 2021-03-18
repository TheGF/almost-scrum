import {
    AccordionButton,
    AccordionIcon,
    AccordionItem,
    AccordionPanel, Box, Button, Input, Switch, Table, Tag, Td, Tr
} from '@chakra-ui/react';
import { React, useState } from "react";
import T from '../../core/T';
import Utils from '../../core/utils';



function S3Exchange(props) {
    const [exchange, setExchange] = useState(props.exchange)
    const { update, status } = props

    const connected = status && status.exchanges[exchange.name]
    const upload = status && status.stat[exchange.name] && Utils.getFriendlySize(status.stat[exchange.name].upload)
    const download = status && status.stat[exchange.name] && Utils.getFriendlySize(status.stat[exchange.name].download)

    function updateField(evt, field, value) {
        exchange[field] = value != undefined ? value : evt.target.value
        const u = { ...exchange }
        setExchange(u)
        update(u)
    }

    return <AccordionItem>
        <h2>
            <AccordionButton>
                <Box flex="1" textAlign="left">
                    S3 - {exchange.name}
                </Box>
                {connected ?
                    <Tag colorScheme="green" title={`U:${upload}-D:${download}`}><T>connected</T></Tag> :
                    <Tag colorScheme="red"><T>disconnected</T></Tag>}
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
                    <Td><label><T>endpoint</T></label></Td>
                    <Td><Input type="text" value={exchange.endpoint} onChange={e => updateField(e, 'endpoint')} /></Td>
                </Tr>
                <Tr>
                    <Td><label><T>bucket</T></label></Td>
                    <Td><Input type="text" value={exchange.bucket} onChange={e => updateField(e, 'bucket')} /></Td>
                </Tr>
                <Tr>
                    <Td><label><T>access key</T></label></Td>
                    <Td><Input type="text" value={exchange.accessKey} onChange={e => updateField(e, 'accessKey')} /></Td>
                </Tr>
                <Tr>
                    <Td><label><T>secret</T></label></Td>
                    <Td><Input type="password" value={exchange.secret} onChange={e => updateField(e, 'secret')} /></Td>
                </Tr>
                <Tr>
                    <Td><label><T>use SSL</T></label></Td>
                    <Td><Switch isChecked={exchange.useSSL} onChange={e => updateField(e, 'useSSL', !exchange.useSSL)} /></Td>
                </Tr>
                <Tr>
                    <Td><label><T>location</T></label></Td>
                    <Td><Input type="text" value={exchange.location} onChange={e => updateField(e, 'location')} /></Td>
                </Tr>
            </Table>

        </AccordionPanel>
    </AccordionItem>
}
export default S3Exchange;