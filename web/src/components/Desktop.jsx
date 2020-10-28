import React, { useState } from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import NavDropdown from 'react-bootstrap/NavDropdown';
import {IoMdLogOut} from 'react-icons/io';
import {GiGrapes} from 'react-icons/gi';
import ProjectsList from './ProjectsList';
import Project from './Project';

function TopBar() {
    function logout() {
        localStorage.removeItem('username');
        localStorage.removeItem('token');
        window.location.assign(window.location.href);
    }


    return  <Navbar bg="dark" variant="dark">
    <Navbar.Brand href="#home"><GiGrapes/> Almost Scrum 0.1</Navbar.Brand>
    <Nav className="ml-auto" >
      <Nav.Link href="#">Git Import</Nav.Link>
      <Nav.Link href="#">Features</Nav.Link>
      <Nav.Link href="#" onClick={logout}>Logout <IoMdLogOut/></Nav.Link>
    </Nav>
    
  </Navbar>
}

function Desktop(props) {
    const [project, setProject] = useState('');

    return <Container fluid>
        <Row>
            <Col md="12">
            <TopBar/>
            </Col>
        </Row>
        <Row>
            <Col md="3">
                <ProjectsList onSelect={setProject} />
            </Col>
            <Col md="9">
                <Project username={props.username} project={project} />
            </Col>
        </Row>
    </Container>;
}

export default Desktop;