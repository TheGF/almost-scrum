import {
    Box, ButtonGroup, HStack, IconButton, Image, Spacer, Tag, useToast, VStack, Wrap, WrapItem,
} from "@chakra-ui/react";
import React, { useContext, useEffect, useState } from 'react';
import { BsTrash, FaRegThumbsUp, MdAddBox, MdFileDownload, MdLibraryAdd, MdLocalGroceryStore } from 'react-icons/all';
import ReactPlayer from "react-player";
import MarkdownEditor from "../core/MarkdownEditor";
import Server from '../server';
import UserContext from '../UserContext';
import ImageViewer from './ImageViewer';
import MakePost from "./MakePost";

function Messages(props) {
    const { project } = useContext(UserContext)
    const [msgs, setMsgs] = useState([])
    const [imageViewerValue, setImageViewerValue] = useState(null)
    const [makePostValue, setMakePostValue] = useState(null)
    const toast = useToast()

    function getMsgs() {
        Server.getChatMsgs(project, 0, 10)
            .then(msgs => setMsgs(msgs || []))
    }
    useEffect(getMsgs, [])


    function TextMsg(props) {
        const { msg } = props

        function deleteChat() {
            Server.deleteChatMsgs(project, msg.id)
                .then(getMsgs)
        }


        function renderTextMsg(text) {
            if (!text || !text.length) return null

            if (text.startsWith("http://") || text.startsWith("https://")) {
                return <HStack width="100%">
                    <ReactPlayer url={text} controls={true} />
                </HStack>
            }
            return <MarkdownEditor readOnly={true} value={text} />
        }

        function renderAttachment(name, mime, idx) {
            function toLibrary() {
                Server.postChatAttachmentAction(project, msg.id, idx, 'make_doc')
                    .then(_ => toast({
                        title: 'Created',
                        description: `The file ${name} has been successfully created in the library`,
                        status: "success",
                        isClosable: true,
                    }))
            }
            function downloadAttachment() {
                downloadLink && downloadLink.click()
            }
            function showImage() {
                setImageViewerValue({
                    project: project,
                    msg: msg,
                    idx: idx
                })
            }


            let downloadLink = null
            const url = Server.getChatAttachmentURL(project, msg.id, idx)

            const buttons = <VStack>
                <a id="link" href="" download={name} ref={l => downloadLink = l} hidden></a>
                <IconButton size="sm" title="Add to Library" onClick={toLibrary}>
                    <MdLibraryAdd />
                </IconButton>
                <IconButton size="sm" title="Download" onClick={downloadAttachment}>
                    <MdFileDownload />
                </IconButton>
            </VStack>

            if (mime.startsWith("audio/")) {
                return <audio key={idx} title={name} controls style={{ height: '24px', width: '100%' }}>
                    <source src={url} type={mime} />
                </audio>
            }
            if (mime.startsWith("image/")) {
                return <WrapItem key={idx}>
                    <HStack width="100%">
                        <Image boxSize="160px" src={url} onClick={showImage} />
                        {buttons}
                    </HStack>
                </WrapItem>
            }
            if (mime.startsWith("video/")) {
                return <HStack key={idx} width="100%">
                    <ReactPlayer
                        width="100%" height="100%"
                        url={url} controls={true} />
                    {buttons}
                </HStack>
            }
        }

        function openMakePost() {
            setMakePostValue(msg)
        }

        const attachmentsUI = msg.names && msg.names.map((n, idx) =>
            renderAttachment(n, msg.mimes[idx], idx))

        const makePostButton = msg.text && msg.text.length > 0 ?
            <IconButton title="make a Post" onClick={openMakePost}><MdAddBox /></IconButton> :
            null

        return <Box key={msg.id} borderWidth={1} borderColor="lightgray" width="100%">
            <VStack width="100%">
                <HStack className="panel2" width="100%">
                    <Tag>{msg.user}</Tag>
                    <Spacer />
                    <ButtonGroup size="sm">
                        <IconButton title="Do you like?"><FaRegThumbsUp /></IconButton>
                        {makePostButton}
                        <IconButton onClick={deleteChat}><BsTrash /></IconButton>
                    </ButtonGroup>
                </HStack>
                {renderTextMsg(msg.text)}
                <Wrap>
                    {attachmentsUI}
                </Wrap>
            </VStack>
        </Box>
    }


    return <VStack spacing="1">
        <MakePost value={makePostValue} onClose={_ => setMakePostValue(null)} />
        <ImageViewer key={!imageViewerValue} value={imageViewerValue} onClose={_ => setImageViewerValue(null)} />
        {msgs.map((m, idx) => <TextMsg key={idx} msg={m} />)}
    </VStack>
}
export default Messages;