import {
    Box, Button, ButtonGroup, Center, CircularProgress, HStack, Image, Spacer, Switch, Table,
    Td, Text, Th, Thead, Tr, useToast, VStack
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import { BiCopy } from 'react-icons/bi';
import Server from '../server';
import UserContext from '../UserContext';

const images = ['banana', 'beer', 'burger', 'chocolate', 'coke', 'cornflakes',
    'cupcake', 'egg', 'fries', 'hotdog', 'juice', 'muffin', 'orange',
    'pasta', 'pizza', 'popcorn', 'steak', 'sushi', 'water', 'wine'
]

function Invite(props) {
    const { project } = useContext(UserContext)
    const { onClose } = props
    const [config, setConfig] = useState(null)
    const [selection, setSelection] = useState([])
    const [invite, setInvite] = useState(null)
    const [keys, setKeys] = useState([])
    const [removeCredentials, setRemoveCredentials] = useState(false)
    const toast = useToast()

    function getConfig() {
        Server.getFedTransport(project)
            .then(setConfig)
    }
    useEffect(getConfig, [])

    function getInvite() {
        const first = Math.floor(Math.random() * images.length)
        let second = Math.floor(Math.random() * images.length)
        while (second == first) second = Math.floor(Math.random() * images.length)

        const keys = [images[first], images[second]].sort()
        Server.postFedShare(project, keys.join(','), selection, removeCredentials)
            .then(invite => {
                setInvite(invite)
                setKeys(keys)
            })
    }


    function chooseExchange(exchange) {
        function change() {
            if (selection.includes(exchange.name)) {
                setSelection(selection.filter(s => s != exchange.name));
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
        function copy(id, description) {
            const copyText = document.getElementById(id);
            const range = document.createRange();
            range.selectNode(copyText);
            window.getSelection().addRange(range);
            // copyText.select();
            // copyText.setSelectionRange(0, 99999); /* For mobile devices */
            document.execCommand("copy");

            //navigator.clipboard.writeText(copyText)
            toast({
                title: `Copied`,
                description: description,
                status: "info",
                duration: 9000,
                isClosable: true,
            })
        }

        function sendByEmail() {
            const msg = document.getElementById('inviteMessage').innerText
            window.location = "mailto:?subject=Invite to Ash together&body=" + encodeURIComponent(msg);
        }

        const keysUI = keys.map(n => <Box borderWidth={1} p={1}>
            <HStack>
                <Image boxSize="32px" src={`/icons/${n}.svg`} ></Image>
                <label>{n}</label>
            </HStack>
        </Box>)

        return <VStack>
            <Box id="inviteMessage">
                <b>Dear ASHer(s)</b>,<br />
            this is an invitation to join the federation.
            <br /><br />
                <a href={`http://localhost:8375?invite=${invite}`}>
                    <font color="blue">Click here</font>
                </a>&nbsp;or copy and paste the following token in the invite dialog<br />
                <pre id="inviteToken" style={{display: 'inline'}}>{multiline(invite)}</pre>
                <BiCopy onClick={_ => copy('inviteToken', '')} style={{display: 'inline'}} />
                <br />
                <HStack>
                    <Text>When requested about Gopher's dinner, choose </Text>
                    {keysUI}
                </HStack>
                <br />
            </Box>
            <ButtonGroup>
                <Button colorScheme="blue" onClick={
                    _ => copy('inviteMessage', 'Now paste in your favorite e-mail client and send')
                }>
                    Copy
                </Button>
                <Button colorScheme="blue" onClick={sendByEmail}>E-mail (text)</Button>
                <Button onClick={_ => setInvite(null)}>Back</Button>
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
                <Button colorScheme="blue" disabled={selection.length == 0} onClick={getInvite}>Get Invite</Button>
                <Button onClick={onClose}>Close</Button>
            </ButtonGroup>
        </VStack>
    }

    return config ? invite ? inviteUI() : createInviteUI() :
        <Center><CircularProgress isIndeterminate color="green.300" size="100px" /></Center>
}
export default Invite;