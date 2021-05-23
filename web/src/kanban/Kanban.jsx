
import { Button, ButtonGroup, Checkbox, HStack, Menu, MenuButton, MenuItem, MenuList, Spacer, Stack, Text } from '@chakra-ui/react';
import Board from '@lourenci/react-kanban';
import '@lourenci/react-kanban/dist/styles.css';
import { React, useContext, useEffect, useState } from 'react';
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';
import MarkdownEditor from '../core/MarkdownEditor';
import useBoardOptions from './BoardOptions';
import Options from './Options';




function Kanban(props) {
  const [board, setBoard] = useState(null)
  const [cardDragEnd, setCardDragEnd] = useState(null)
  const [redraw, setRedraw] = useState(false)


  function getColumnsFromRefs(columnTitles, refs, getColumnTitle) {
    function addCard(columns, ref) {
      const title = getColumnTitle(ref)
      for (const c of columns) {
        if (c.title == title) {
          const description = ref.task.description.substring(0, 64) + ref.task.description > 64 ? '...' : ''
          const owner = ref.task.properties && ref.task.properties['Owner']
          const cardTitle = <HStack><Text>{ref.name}</Text><Spacer /><Text fontSize="sm">{owner}</Text></HStack>
          c.cards.push({
            id: ref.name,
            title: cardTitle,
            description: <MarkdownEditor value={description} height={200} readOnly={true} />,
            ref: ref,
          })
        }
      }
    }

    const columns = columnTitles.map((title, idx) => ({
      id: idx + 1,
      title: title,
      cards: [],
    }))
    for (const ref of refs) {
      addCard(columns, ref)
    }
   return columns
  }

  function updateBoard(columnTitles, refs, getColumnTitle, cardDragEnd) {
    const columns = getColumnsFromRefs(columnTitles, refs, getColumnTitle)
    setBoard(<Board key={new Date()} initialBoard={{columns: columns}} onCardDragEnd={cardDragEnd} />)
//    setBoard({columns: columns})
//    setCardDragEnd(cardDragEnd)
  }

  return <Stack className="yellowPanel" w="100%" direction="column" spacing={2} p={2}>
    <Options updateBoard={updateBoard}/>
    {board}
  </Stack>

}

export default Kanban
