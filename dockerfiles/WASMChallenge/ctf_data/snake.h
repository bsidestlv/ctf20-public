extern uint16_t get_random(uint16_t min, uint16_t max);
extern void draw_at(uint16_t x, uint16_t y, uint32_t color);
extern void console_log(uint32_t);
extern void draw_game_lost();
extern void draw_clear_board();
void game_init();

#define CELL_SIZE 16
#define WIDTH 576
#define HEIGHT 512
#define EMPTY_COORD -1
#define BODY 0
#define BOARD_TILE 1
#define FRUIT_1 2
#define FRUIT_2 3

#ifdef FLAG_BYTES
const unsigned char const_arr[] = {FLAG_BYTES};
#else
const unsigned char const_arr[] = {0x31,0x11,0x3a,0xd,0x1,0x16,0x27,0x18,0x1a,0x64,0x2,0x2,0x2,0x4b,0x2c,0x66,0x42,0x17,0xb,0x2,0x32,0x36,0x5c,0x6a,0x30,0x1,0x2,0x15,0x26,0x2f,0x3f,0x3c,0x46,0x50,0x16,0x0,0x40,0x57,0x3b,0x3d,0x1b,0x26,0x2b,0x1c,0x5b,0x6c,0x33,0x5f,0x7,0x46,0x1c,0xb,0x1,0x19};
#endif

const char initial_key = 's';

typedef enum {
    KEY_UNCHANGED = 0,
    KEY_REPLAY = 1,
    KEY_RIGHT = 3,
    KEY_DOWN = 5,
    KEY_UP = 10,
    KEY_LEFT = 12,
    KEY_SPACE = 15
} input_t;

typedef enum {
    STOPPED,
    STARTED,
    LOST,
    ATE,
    NORMAL,
} state_t;

typedef struct {
    uint16_t x;
    uint16_t y;
} pos_t;

typedef struct {
    pos_t position;
    uint32_t color;
    uint16_t len;
    input_t input;
    state_t state;
} player_t;

typedef struct {
    pos_t position;
    uint32_t color;
    uint32_t iteration;
} fruit_t;

typedef struct {
    player_t player;
    fruit_t fruit;
    uint32_t time;
    state_t state;
    pos_t border;
} game_t;


extern input_t get_input();