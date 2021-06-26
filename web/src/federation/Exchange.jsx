import {
    AccordionButton,
    AccordionIcon,
    AccordionItem,
    AccordionPanel, Box, Button, Input, Switch, Table, Tag, Td, Tr
} from '@chakra-ui/react';
import { React, useState, useContext, useEffect } from "react";
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import models from './models'

let delayedUpdate = null;

function Exchange(props) {
    const { project } = useContext(UserContext)
    const { id, status } = props;
    const [config, setConfig] = useState(null)
    const modelId = id.split('/')[0]
    const model = models[modelId]

    function getExchange() {
        Server.getFedTransportExchange(project, id)
            .then(setConfig)
    }
    useEffect(getExchange, []);

    function updateField(evt, name, value) {
        const c = {...config}
        c[name] = value != undefined ? value : evt.target.value
        setConfig(c)
        
        clearTimeout(delayedUpdate)
        delayedUpdate = setTimeout(_=>Server.putFedTransportExchange(project, id, c), 2000)
    }

    function getField(name) {
        switch (typeof model[name]) {
            case 'string': return name in ['password', 'secret'] ?
                <Input type="password" value={config[name]} onChange={e => updateField(e, name)} /> :
                <Input type="text" value={config[name]} onChange={e => updateField(e, name)} />
            case 'number':
                return <Input type="number" value={config[name]} onChange={e => updateField(e, name)} />
            case 'boolean':
                return <Switch isChecked={config[name]} onChange={e => updateField(e, name, !config[name])} />
        }
    }

    const rows = config && Object.keys(model).map(k =>  <Tr>
        <Td><label><T>{k}</T></label></Td>
        <Td>{getField(k)}</Td>
    </Tr>)

    return <AccordionItem>
        <h2>
            <AccordionButton>
                <Box flex="1" textAlign="left">
                    {config && config.name}
                </Box>
                {status == '' ?
                    <Tag colorScheme="green"><T>connected</T></Tag> :
                    <Tag colorScheme="red">{status}</Tag>}
                <AccordionIcon />
            </AccordionButton>
        </h2>
        <AccordionPanel w="100%">
            <Button w="100%" onClick={_ => props.update(null)}>Delete</Button>
            <Table size="sm" padding="0">
                {rows}
            </Table>

        </AccordionPanel>
    </AccordionItem>
}
export default Exchange;