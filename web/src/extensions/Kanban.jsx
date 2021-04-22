
import { React, useContext, useEffect, useState } from 'react';
import Server from '../server';
import UserContext from '../UserContext';
import { Box, Button, ButtonGroup, Text, Stack, HStack, Spacer, Select, VStack } from '@chakra-ui/react';

import Board from '@lourenci/react-kanban'
import '@lourenci/react-kanban/dist/styles.css'


function Kanban(props) {
    const board = {
        columns: [
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
      }
      
      return <VStack className="panel1">
          <Text>Status, Owner, Boards</Text>
      <Board initialBoard={board} />
      </VStack>
      
}

export default Kanban