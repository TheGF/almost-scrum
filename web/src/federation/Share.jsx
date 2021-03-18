import {
    Accordion, Button, ButtonGroup, Input, Spacer, Switch, Table, Td, Textarea, Tr, VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import S3Exchange from './exchanges/S3Exchange';
import FTPExchange from './exchanges/FTPExchange';
import WebDAVExchange from './exchanges/WebDAVExchange';


function Share(props) {
    const { project, onClose } = useContext(UserContext)
    const [config, setConfig] = useState(null)
    const [selection, setSelection] = useState({})
    const [tokens, setTokens] = useState(null)

    function getConfig() {
        Server.getFedConfig(project)
            .then(setConfig)
    }
    useEffect(getConfig, [])

    function getTokens() {
        Server.postFedShare(project, Object.keys(selection), false)
            .then(setTokens)
    }


    function chooseExchange(exchange) {
        function change() {
            setSelection({
                ...selection,
                exchange: !selection[exchange]
            })
        }

        return <Tr>
            <Td>{exchange.name}</Td>
            <Td><Switch isChecked={selection[exchange.name]}
                onChange={change} /></Td>
        </Tr>

    }

    let exchangesUI = config ? [
        ...config.s3,
        ...config.webDAV,
        ...config.ftp,
    ].map(exchange => chooseExchange(exchange)) : []

    function multiline(s) {
        let r = ''
        for (let i = 0; i < s.length; i++) {
            if (i % 64 == 0) r += '\n\t'
            r += s[i]
        }
        return r
    }

    return config ? tokens ?
        <Textarea rows={20} value={`
            Dear Asher(s),
            please use the following token to access the federation
            ${multiline(tokens.token)} 

            When requested for the key, insert the following.
            ${tokens.key} 

        `}/>
        : <VStack>
            <Table size="sm" padding="0" spacing="0">
                {exchangesUI}
            </Table>
            <Spacer minHeight="1em" />
            <ButtonGroup>
                <Button colorScheme="blue" onClick={getTokens}>Get Tokens</Button>
                <Button onClick={onClose}>Close</Button>
            </ButtonGroup>
        </VStack> : <span>Loading</span>
}
export default Share;