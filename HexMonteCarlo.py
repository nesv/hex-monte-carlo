#!/usr/bin/python

import sys, numpy, random, copy, time


# Run Monte-Carlo simulations of Hex.

# TODO:
# -- Implement Gale's algorithm for determining who wins a filled board
# -- Generate random Hex boards (Ramsey, Standard, and Symmetric)

# the standard board structure:
#    boards are n x n lists of lists, with each entry being '1'st player or '2'nd player
#    1st player connects vertically, 2nd player horizontally

# the board convention

#    22222
#  1 o-o-o 1
#  1 |/|/| 1
#  1 o-o-o 1
#  1 |/|/| 1
#  1 o-o-o 1
#    22222

#  the indexing convention

# [y][x] = [y,x]

#   ----------> 
# | 0,0 0,1 0,2
# | 1,0 1,1 1,2
# v 2,0 2,1 2,2
#               

# example boards

# 1st player wins
example_first_a  = [[1,1,1],
                    [1,1,1],
                    [1,1,1]];

# 2nd player wins
example_second_a  = [[2,2,2],
                     [2,2,2],
                     [2,2,2]];
 
def extend(board):
    """ Mark the edges of the board with stones denoting who needs to connect.
        This creates the border of stones used in Gale's Algorithm. 

                    21111
        a b c       2abc2
        e f g ----> 2efg2
        h i j       2hij2
                    11112
        
        This convention dictates: 
            -- 2nd plays horizontally
            -- 1st plays vertically
    """

    sideLength = len(board);

    extended_board = [];

    first_row = [2] + ( [1] * (sideLength + 1) )
    extended_board.append(first_row);

    for row in board:
        middle_row = [2] + row + [2];
        extended_board.append(middle_row);

    last_row = ([1] * (sideLength + 1) ) + [2];
    extended_board.append(last_row);

    return extended_board;
 
def printBoard(board):
    """ Print a nice grid layout of the board """
    for row in board:
        print "  ".join(map(str,row));

def onBoard(v,n):
    """ Given a vector, determine if it is on the board. """
    if ((-1 <= v[0] <= n + 1) and (0 <= v[1] <= n + 1)): 
        return True;
    else:
        return False;
    

def Gale(board):
    """ Given a filled board, return the winner. """
    xBoard = extend(board);
    sideLength = len(board[0]);

    # the four vertices in Gale's Algorithm
    # 
    #           v1  [-1,1]
    #          / |
    #         /  |
    #        /   |
    #       /    |
    #[0,0] v2-|-v3 [0,1]
    #      |  v  /
    #      |    /
    #      |   /
    #      |  /
    #       v4 [1,0]
    #
    # note: left and right are revered, since we head along the arrow

    v1 = numpy.array([-1,1]);
    v2 = numpy.array([0,0]);
    v3 = numpy.array([0,1]);
    v4 = numpy.array([1,0]);

    while (onBoard(v4, sideLength)):
        v1_new = numpy.array;
        v2_new = numpy.array;
        v3_new = numpy.array;
        v4_new = numpy.array;

        # PRINT the edge just passed
        #print("\t {0}-{1}".format(str(v2), str(v3)))

        if xBoard[v4[0]][v4[1]] == xBoard[v3[0]][v3[1]]:
            # go left 
            v1_new = v3;
            v2_new = v2;
            v3_new = v4;
            v4_new = v2 + (v2 - v1);
        elif xBoard[v4[0]][v4[1]] == xBoard[v2[0]][v2[1]]:
            # go right
            v1_new = v2;
            v2_new = v4;
            v3_new = v3;
            v4_new = v3 + (v3 - v1);
        
        v1,v2,v3,v4=v1_new,v2_new,v3_new,v4_new;
  
    # [n,-1] = bottom left
    # if v4 exits the bottom left corner, then `1'st won
    if v4[1] == -1:
        return 1;
    # [0,n] = top right
    # if v4 exits the top right corner, then `2'nd won
    elif v4[0] == 0:
        return 2;
    else:
       raise ValueError('The Gale algorithm broke: not [* -1] or [0 *].')

def BernoulliRandomBoard(sidelength):
    """ Produces a filled board by tossing a fair {1,2}-coin for each cell. """

    board = [[random.choice([1,2]) for x in range(0, sidelength)] for x in range(0,sidelength)];

    return(board);

### Try out Gale on random boards ###

#for i in range(1,1000):
#    board = extend(BernoulliRandomBoard(10))
#    #printBoard(board)
#    print(Gale(board))

#### Simulator ####

N = 1000;

halfBoardSize = 3;
BoardSize = 2*halfBoardSize + 1
criticalCount = numpy.zeros((BoardSize, BoardSize));

print "Simulating {0}x{0} with N={1} trials".format(BoardSize,N);

for x in range(0,BoardSize):
    for y in range(0, BoardSize):

        for i in range(1,N+1):
            board = BernoulliRandomBoard(BoardSize);

            # we need to copy the array, otherwise python just points to board
            board1 = copy.deepcopy(board);
            board1[y][x] = 1;

            board2 = copy.deepcopy(board);
            board2[y][x] = 2;
       
            if (Gale(board1) != Gale(board2)):
                criticalCount[y][x] += 1;   

print criticalCount;
