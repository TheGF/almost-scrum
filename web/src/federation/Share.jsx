import {
    Box, Button, ButtonGroup, Center, CircularProgress, HStack, Spacer, Switch, Table,
    Td, Th, Thead, Tr, useToast, VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import Server from '../server';
import UserContext from '../UserContext';

function Share(props) {
    const { project } = useContext(UserContext)
    const { onClose } = props
    const [config, setConfig] = useState(null)
    const [selection, setSelection] = useState([])
    const [invite, setInvite] = useState(null)
    const [removeCredentials, setRemoveCredentials] = useState(false)
    const toast = useToast()

    function getConfig() {
        Server.getFedConfig(project)
            .then(setConfig)
    }
    useEffect(getConfig, [])

    function getInvite() {
        Server.postFedShare(project, selection, removeCredentials)
            .then(setInvite)
    }


    function chooseExchange(exchange) {
        function change() {
            if (selection.includes(exchange.name)) {
                setSelection(selection.filter(s=>s!=exchange.name));
            } else {
                setSelection([...selection, exchange.name])
            }
        }

        return <Tr>
            <Td>{exchange.name}</Td>
            <Td><Switch isChecked={selection.includes(exchange.name)}
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
            if (i % 100 == 0) r += '\n'
            r += s[i]
        }
        return r
    }

    function inviteUI() {
        function copy() {
            const copyText = document.getElementById("inviteMessage");
            const range = document.createRange();
            range.selectNode(copyText);
            window.getSelection().addRange(range);
            // copyText.select();
            // copyText.setSelectionRange(0, 99999); /* For mobile devices */
            document.execCommand("copy");

            //navigator.clipboard.writeText(copyText)
            toast({
                title: `Copied`,
                description: 'Now paste in your favorite e-mail client and send',
                status: "info",
                duration: 9000,
                isClosable: true,
            })
          } 

        function sendByEmail() {
            const msg = document.getElementById('inviteMessage').innerText
            window.location="mailto:?subject=Invite to Ash together&body="+encodeURIComponent(msg);
        }

        return <VStack>
            <Box id="inviteMessage">
                <b>Dear Asher(s)</b>,<br />
            this is an invitation to join the federation.
            <br /><br />
                <a href={`http://localhost:8375?invite&token=${invite.token}`}>
                    <font color="blue">Click here</font>
                </a>&nbsp;or copy and paste the following token in the invite dialog<br />
                <pre id="inviteToken">{multiline(invite.token)}</pre>
                <br />When requested for the decryption key, insert <b>{invite.key}</b>
            </Box>
            <ButtonGroup>
                <Button colorScheme="blue" onClick={copy}>Copy</Button>
                <Button colorScheme="blue" onClick={sendByEmail}>E-mail (text)</Button>
                <Button onClick={_=>setInvite(null)}>Back</Button>
                <Button onClick={onClose}>Close</Button>
            </ButtonGroup>
        </VStack>
    }

    function createInviteUI() {
        return <VStack>
            <HStack>
                <span>Remove credentials (password and secrets) from the invite</span>
                <Switch isChecked={removeCredentials}
                    onChange={_ => setRemoveCredentials(!removeCredentials)} />
            </HStack>

            <br />
            <Table size="sm" padding="0" spacing="0">
                <Thead>
                    <Tr>
                        <Th>Exchange</Th>
                        <Th>Add to invite</Th>
                    </Tr>
                </Thead>
                {exchangesUI}
            </Table>
            <Spacer minHeight="1em" />
            <ButtonGroup>
                <Button colorScheme="blue" onClick={getInvite}>Get Invite</Button>
                <Button onClick={onClose}>Close</Button>
            </ButtonGroup>
        </VStack>
    }

    return config ? invite ? inviteUI() : createInviteUI() :
        <Center><CircularProgress isIndeterminate color="green.300" size="100px" /></Center>
}
export default Share;