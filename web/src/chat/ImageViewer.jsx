import {
    Box, ButtonGroup, HStack, IconButton, Image, Spacer, Tag, VStack, Wrap, WrapItem,
} from "@chakra-ui/react";
import React, { useContext, useEffect, useState } from 'react';
import { BsTrash, FaRegThumbsUp, MdFileDownload, MdLibraryAdd, MdLocalGroceryStore, MdNavigateBefore, MdNavigateNext } from 'react-icons/all';
import ReactPlayer from "react-player";
import MarkdownEditor from "../core/MarkdownEditor";
import Server from '../server';
import UserContext from '../UserContext';

import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
} from "@chakra-ui/react"


function ImageViewer(props) {
    const { value, onClose } = props;
    const [idx, setIdx] = useState(value ? value.idx : null)

    const project = value && value.project;
    const msg = value && value.msg
    const indexes = msg && msg.mimes && msg.mimes.map((m, idx) => m.startsWith('image/') ? idx : -1).filter(i => i >= 0)
    const first = indexes && indexes[0]
    const last = indexes && indexes[indexes.length - 1]
    const url = msg && Server.getChatAttachmentURL(project, msg.id, idx)
    const name = msg && msg.names[idx]

    function moveTo(step) {
        const i = indexes.indexOf(idx)
        setIdx(indexes[i + step])
    }

    return <Modal isOpen={value != null} onClose={onClose} size="2xl">
        <ModalOverlay />
        <ModalContent>
            <ModalHeader>{name}</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
                <HStack width="100%">
                    <IconButton disabled={idx == first} onClick={_ => moveTo(-1)}>
                        <MdNavigateBefore />
                    </IconButton>
                    <Image src={url} width="90%" />
                    <IconButton disabled={idx == last} onClick={_ => moveTo(1)}>
                        <MdNavigateNext />
                    </IconButton>
                </HStack>
            </ModalBody>
            <ModalFooter>
            </ModalFooter>
        </ModalContent>
    </Modal>
}

export default ImageViewer