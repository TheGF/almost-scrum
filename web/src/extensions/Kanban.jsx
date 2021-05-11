
import { React, useContext, useEffect, useState } from 'react';
import Server from '../server';
import UserContext from '../UserContext';
import { Box, Button, ButtonGroup, Text, Stack, HStack, Spacer, Select, VStack, Checkbox } from '@chakra-ui/react';

import Board from '@lourenci/react-kanban'
import '@lourenci/react-kanban/dist/styles.css'
import { getOptions } from 'showdown';
import T from '../core/T';


const columns = [
  {
    id: 1,
    title: 'Backlog',
    cards: [
      {
        id: 1,
        title: 'Add card',
        description: 'Add capability to add a card in a column'
      },
    ]
  },
  {
    id: 2,
    title: 'Doing',
    cards: [
      {
        id: 2,
        title: 'Drag-n-drop support',
        description: 'Move a card between the columns'
      },
    ]
  }
]

function Kanban(props) {
  const { info, project, reload } = useContext(UserContext)
  const [view, setView] = useState('boards')
  const [boards, setBoards] = useState([])
  const [columns, setColumns] = useState([])

  function listBoards() {
    Server.listBoards(project)
      .then(setBoards)
  }
  useEffect(_=>listBoards(), [])

  function BoardOptions(props) {
    const boardsUI = boards.map(b => <Checkbox>{b}</Checkbox>)
    const [selectedBoards, setSelectedBoards] = useState({})

    return <HStack spacing={3}>{boardsUI}</HStack>
  }

  function getOptions() {
    return <BoardOptions/>
  }

  const board = {
    columns: columns,
  }

  const viewsUI = ['boards', 'people', 'status', 'property'].map(v =>
    <Button isActive={view == v} onClick={_ => setView(v)}><T>{v}</T></Button>
  )

  return <Stack className="yellowPanel" w="100%" direction="column" spacing={2} p={2}>
    <HStack>
      <ButtonGroup size="sm" >{viewsUI}</ButtonGroup>
    </HStack>
    {getOptions()}

    <Board initialBoard={board} />
  </Stack>

}

export default Kanban