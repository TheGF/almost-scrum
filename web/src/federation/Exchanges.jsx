import {
    Accordion, Button, ButtonGroup, FormLabel, Input, Spacer, Table, Td, Tr, VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import Exchange from './Exchange';
import FTPExchange from './exchanges/FTPExchange';
import S3Exchange from './exchanges/S3Exchange';
import USBExchange from './exchanges/USBExchange';
import WebDAVExchange from './exchanges/WebDAVExchange';
import models from './models';


function Exchanges(props) {
    const { project } = useContext(UserContext)
    const { onClose } = props
    const [transport, setTransport] = useState(null)

    function init() {
        Server.getFedTransport(project)
            .then(setTransport)
    }
    useEffect(init, [])

    function addExchange(modelId) {
        const model = models[modelId]
        Server.putFedTransportExchange(project, `${modelId}/${new Date().getTime()}`, model)
            .then(init)
    }

    let exchangesUI = transport && Object.keys(transport)
                            .sort()
                            .map(k => <Exchange id={k} status={transport[k]} />) || []

    const addButtons = Object.keys(models).map(k => <Button onClick={_ => addExchange(k)}><T>{`add ${k}`}</T></Button>)


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
            {addButtons}
            <Spacer minW="2em" />
            <Button onClick={onClose}>Close</Button>
        </ButtonGroup>
        <Accordion w="100%" allowToggle>
            {exchangesUI}
        </Accordion>
    </VStack> : <span>Loading</span>
}
export default Exchanges;