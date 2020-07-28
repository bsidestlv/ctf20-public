#include <stdint.h>
#include <emscripten.h>
#include "snake.h"

game_t game = {0};
pos_t snake_coord_ptr[WIDTH][HEIGHT];
pos_t head;
pos_t tail;
char key;

uint32_t get_next_fruit()
{
    int index = game.fruit.iteration;
    if (!(index % 8))
        key = key ^ const_arr[index / 8];
    if (key & (1 << (7 - (index % 8))))
        return FRUIT_1;
    return FRUIT_2;
}

void draw()
{
    // draw players head
    draw_at(game.player.position.x, game.player.position.y, BODY);

    // draw fruit
    draw_at(game.fruit.position.x, game.fruit.position.y, game.fruit.color);

}

void move_player(player_t * player)
{
    // TODO: check borders
    pos_t temp_head = head;
    switch (player->input)
    {
        case KEY_DOWN:
           player->position.y += CELL_SIZE;
            break;
        case KEY_UP:
           player->position.y -= CELL_SIZE;
            break;
        case KEY_LEFT:
           player->position.x -= CELL_SIZE;
            break;
        case KEY_RIGHT:
           player->position.x += CELL_SIZE;
            break;
        case KEY_SPACE:
            game.state = STOPPED;
        default:
            break;
    }
    head = (pos_t){player->position.x, player->position.y};
    snake_coord_ptr[temp_head.x][temp_head.y] = head;
}

void remove_tail()
{
    draw_at(tail.x, tail.y, BOARD_TILE);
    pos_t tail_temp = snake_coord_ptr[tail.x][tail.y];
    snake_coord_ptr[tail.x][tail.y] = (pos_t){EMPTY_COORD, EMPTY_COORD};
    tail = tail_temp;
}

uint8_t check_border_collision()
{
    return  0 > game.player.position.x ||
            game.border.x < game.player.position.x ||
            0 > game.player.position.y || 
            game.border.y < game.player.position.y;
}

uint8_t check_self_collision()
{
    return (snake_coord_ptr[head.x][head.y].x != (uint16_t) EMPTY_COORD);
}

uint8_t is_overlapping_obejcts(pos_t * object1, pos_t * object2)
{
    return object1->x == object2->x && object1->y == object2->y;
}

uint16_t getMultCellSize(uint16_t num)
{
    return (num / CELL_SIZE) * CELL_SIZE;
}

uint8_t update()
{
    game.player.state = NORMAL;
    move_player(&game.player);
    if (check_border_collision() || check_self_collision())
    {
        game.state = LOST;
    }

    if (is_overlapping_obejcts(&game.player.position, &game.fruit.position))
    {
        game.player.state = ATE;
        game.player.len++;
        game.fruit.iteration = (game.fruit.iteration + 1) % sizeof const_arr;
        game.fruit.color = get_next_fruit();
        game.fruit.position.x = getMultCellSize(get_random(CELL_SIZE, game.border.x - CELL_SIZE));
        game.fruit.position.y = getMultCellSize(get_random(CELL_SIZE, game.border.y - CELL_SIZE));
    }

    return 0;
}

EMSCRIPTEN_KEEPALIVE 
void game_iter()
{
    input_t inp = get_input();
    if (KEY_UNCHANGED != inp)
    {
        if (STOPPED == game.state)
        {
            game.state = STARTED;
        }

        if (KEY_REPLAY == inp){
            if (LOST == game.state){
                game.state = STOPPED;
                draw_clear_board();
                game_init();
            }
        }
        // ignore opposite directions
        if ((game.player.input & inp))
            game.player.input = inp;
    }
    
    if (STARTED == game.state)
    {
        update();
    }
    if (STARTED == game.state)
    {
        draw();
        if (ATE != game.player.state){
            remove_tail();
        }

    } 
    if (LOST == game.state){
        draw_game_lost();
    }
}
void init_snake_array(x,y)
{
    int i,j;
    //no memset
    for(i = 0; i< WIDTH; i += CELL_SIZE){
        for (j = 0; j<HEIGHT; j+= CELL_SIZE){
            snake_coord_ptr[i][j] = (pos_t){EMPTY_COORD, EMPTY_COORD};
        }
    }
    //update pointers
    head = (pos_t){x,y};
    tail = (pos_t){x + (2 * CELL_SIZE), y};
    //draw
    draw_at(x, y, BODY);
    draw_at(x + CELL_SIZE, y, BODY);
    draw_at(x + (2 * CELL_SIZE), y, BODY);

    // draw fruit
    draw_at(game.fruit.position.x, game.fruit.position.y, game.fruit.color);

    //update array
    snake_coord_ptr[x + CELL_SIZE][y] = (pos_t)head;
    snake_coord_ptr[x + (2 * CELL_SIZE)][y] = (pos_t){x + CELL_SIZE, y};

}


EMSCRIPTEN_KEEPALIVE 
void game_init()
{
    game.player.position.x = getMultCellSize(get_random(CELL_SIZE * 5, WIDTH - (CELL_SIZE *5)));
    game.player.position.y = getMultCellSize(get_random(5, HEIGHT - 5));
    game.border.x = WIDTH;
    game.border.y = HEIGHT;

    game.fruit.position.x = getMultCellSize(get_random(CELL_SIZE, WIDTH - CELL_SIZE));
    game.fruit.position.y = getMultCellSize(get_random(CELL_SIZE, HEIGHT - CELL_SIZE));

    key = initial_key;
    game.fruit.iteration = 0;
    game.fruit.color = get_next_fruit();

    game.player.input = KEY_LEFT;
    game.player.state = NORMAL;
    
    init_snake_array(game.player.position.x,game.player.position.y);
}