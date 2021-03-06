import React, { useEffect, useState } from 'react';
import { Modal } from 'reactstrap';

import axios from 'axios';
import { render } from 'react-dom';

import { Button, ButtonGroup, HStack, Input, Stack, AspectRatio, Center, Square } from "@chakra-ui/react"
import {
    Table,
    Thead,
    Tbody,
    Tfoot,
    Tr,
    Th,
    Td,
    TableCaption,

    NumberInput,
    NumberInputField,
    NumberInputStepper,
    NumberIncrementStepper,
    NumberDecrementStepper,

} from "@chakra-ui/react"

function AddBook(props) {

    const [id, setId] = useState()
    const [title, setTitle] = useState()
    const [isbnNo, setISBNo] = useState()
    const [series, setSeries] = useState()

    function handleAddBook(e) {
        e.preventDefault()
        axios.post(`http://localhost:3000/api/books`, {
            "id": parseInt(id),
            "title": title,
            "ISBNno": isbnNo,
            "Author": "",
            "Series": series,
            "Genre": "",
            "Rating": 3.8,

        }).then(r => console.log("msg: ", r))
            .catch(e => console.log("err: ", e))
        console.log("add", id, title, isbnNo, series)

    }



    return (
        <AspectRatio maxW="560px" ratio={5 / 3}>
            <Center w="460px" h="460px" >
                <form display="block">

                    <Stack spacing={3}>

                        {/* <Input variant="outline" placeholder="ID"   type="int" value={id} onChange={e=>setId(e.target.value)} name="id"/> */}
                        <NumberInput >
                            <NumberInputField variant="outline" placeholder="ID" type="int" value={id} onChange={e => setId(e.target.value)} name="id"/>
                            <NumberInputStepper>
                                <NumberIncrementStepper />
                                <NumberDecrementStepper />
                            </NumberInputStepper>
                        </NumberInput>


                        <Input variant="outline" placeholder="TITLE" type="text" value={title} onChange={e => setTitle(e.target.value)} name="title" />

                        <Input variant="outline" placeholder="ISBNNo" value={isbnNo} onChange={e => setISBNo(e.target.value)} name="isbno" />

                        <Input variant="outline" placeholder="SERIES" onChange={e => setSeries(e.target.value)} name="series" />


                        <Button onClick={handleAddBook} iconSpacing="-1" size="sm" colorScheme="teal">Add</Button>

                    </Stack>
                </form>
            </Center>
        </AspectRatio>

    )
}





function Book(props) {
    const book = props.book

    function handleDelete() {
        console.log(book)

        axios.delete(`http://localhost:3000/api/books/${book.id}`)
            .then((response) => {
                console.log("success")

            })
            .catch(error => {
                console.log(error);
                // this.setState({...this.state, isFetching: false});
            })

    }
    return (
        <tr key={book.id}>
            <td>{book.id}</td>
            <td>{book.title}</td>
            <td>{book.ISBNno}</td>
            <td>{book.Author}</td>
            <td>{book.Series}</td>
            <td>{book.Genre}</td>
            <td>
                <HStack spacing="24px">
                    <Button size="sm" colorScheme="teal">UPDATE</Button>
                    <Button onClick={handleDelete} iconSpacing="-1" size="sm" colorScheme="red">DELETE</Button>
                </HStack>
            </td>
        </tr>

    )

}

function BookList() {
    const [launchedBooks, setLaunchedBooks] = useState([]);

    // state = [newbookdata= {
    //   id: "",
    //   title: "",
    // }] 


    useEffect(() => {
        axios.get('http://localhost:3000/api/books')
            .then((response) => {
                setLaunchedBooks(response.data)
            })

            .catch(error => {
                console.log(error);
                // this.setState({...this.state, isFetching: false});
            })

    }, [])

    let books = launchedBooks.map((book) => (<Book book={book} />));

    return (
        <div className="App container">


            <Table colorScheme="purple" variant="simple" border="1px solid #605D5D"   >
                <Thead>
                    <Tr>
                        <Th>Id</Th>
                        <Th>Title</Th>
                        <Th>ISBNNo</Th>
                        <Th>Author</Th>
                        <Th>Series</Th>
                        <Th>Genre</Th>
                        <Th>Actions</Th>
                    </Tr>
                </Thead>

                <Tbody>
                    {books}
                </Tbody>

            </Table>


        </div>
    );

}

export default function Home() {
    return (
        <>
            <BookList />
            <AddBook />
        </>
    )
}