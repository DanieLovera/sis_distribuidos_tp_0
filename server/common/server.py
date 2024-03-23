import socket
import logging
from common.signal_manager import SignalManager

class Server:
    def __init__(self, port, listen_backlog):
        signal_manager = SignalManager()
        signal_manager.add_handler(signal_manager.SIGTERM, self._sigterm_handler)
        # Initialize server socket
        self._client_socket = None
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while True:
            self._client_socket = self.__accept_new_connection()
            # Client socket is None when server is shutting down
            if (self._client_socket is None): break
            self.__handle_client_connection()

    def __handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = self._client_socket.recv(1024).rstrip().decode('utf-8')
            addr = self._client_socket.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            self._client_socket.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            self._release_client_socket()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        try:
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except OSError as e:
            if e.errno == 9:
                logging.info(f'action: accept_connections | result: fail | error: Closing server')
                return None
            else:
                raise e

    def _release_resources(self):
        self._release_client_socket()
        self._release_server_socket()

    def _release_server_socket(self):
        if self._server_socket:
            try:
                self._server_socket.shutdown(socket.SHUT_RDWR)
            except OSError as e:
                pass
            finally:
                self._server_socket.close()
                self._server_socket = None

    def _release_client_socket(self):
        if self._client_socket:
            try:
                self._client_socket.shutdown(socket.SHUT_RDWR)
            except OSError as e:
                pass
            finally:
                self._client_socket.close()
                self._client_socket = None

    def _sigterm_handler(self, signum, frame):
        logging.info("action: sigterm_handler | result: in_progress")
        self._release_resources()
        logging.info("action: sigterm_handler | result: success")