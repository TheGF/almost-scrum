import {
    AccordionButton,
    AccordionIcon,
    AccordionItem,
    AccordionPanel, Box, Button, Tag, Text
} from '@chakra-ui/react';
import { React, useState } from "react";
import T from '../../core/T';



function USBExchange(props) {
    const [exchange, setExchange] = useState(props.exchange)
    const { update, status } = props
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
                    USB Media 
                </Box>
                {connectErr == null ?
                    <Tag colorScheme="green" title={`U:${upload}-D:${download}`}><T>connected</T></Tag> :
                    <Tag colorScheme="red">{connectErr && connectErr.Msg ? connectErr.Msg : <T>disconnected</T> }</Tag>}
                <AccordionIcon />
            </AccordionButton>
        </h2>
        <AccordionPanel w="100%">
            <Text>
                Create a folder <i>AlmostScrum</i> on your USB media to enable the 
                synchronization. Then reload the page.
            </Text>
            <Button w="100%" onClick={_=>props.update(null)}>Delete</Button>
        </AccordionPanel>
    </AccordionItem>
}
export default USBExchange;