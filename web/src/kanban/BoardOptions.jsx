import { Button, ButtonGroup, Checkbox, HStack, Menu, MenuButton, MenuItem, MenuList, Spacer, Stack, Text } from '@chakra-ui/react';
import Board from '@lourenci/react-kanban';
import '@lourenci/react-kanban/dist/styles.css';
import { React, useContext, useEffect, useState } from 'react';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import MarkdownEditor from '../core/MarkdownEditor';



function BoardOptions(props) {
    const { project } = useContext(UserContext)
    const { updateBoard } = props
    const [boards, setBoards] = useState([])
    const [selectedBoards, setSelectedBoards] = useState([])

    function listBoards() {
        Server.listBoards(project)
            .then(setBoards)
    }
    useEffect(_ => listBoards(), [])

    function boardOptionsDragEnd(board, card, source, destination) {
        if (board) {
            const s = board.columns[source.fromColumnId - 1]
            const t = board.columns[destination.toColumnId - 1]
            Server.moveTask(project, s.title, card.id, t.title)
        }
    }
    
    function checkBoard(b) {
        const boards = selectedBoards.includes(b) ?
            selectedBoards.filter(s => s != b) :
            [...selectedBoards, b]
        Server.postQueryTasks(project, { select: { description: true, properties: true }, whereBoardIs: boards })
            .then(refs => {
                refs = refs || []
                updateBoard(boards, refs, r => r.board, boardOptionsDragEnd)
            })
        setSelectedBoards(boards)

    }

    const boardsUI = boards.map(b => <Checkbox key={b} onChange={_ => checkBoard(b)} isChecked={selectedBoards.includes(b)}>
        {b}
    </Checkbox>)

    return <HStack spacing={3}>{boardsUI}</HStack>
}

export default BoardOptions