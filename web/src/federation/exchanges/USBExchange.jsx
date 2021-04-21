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
                {props.connected ?
                    <Tag colorScheme="green"><T>connected</T></Tag> :
                    <Tag colorScheme="red"><T>disconnected</T></Tag>}
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