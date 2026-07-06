//! TigerCasino Signature Games Module
//! 
//! Implements provably fair versions of Crash, Mines, Plinko, and other signature games

pub mod crash;
pub mod mines;
pub mod plinko;
pub mod dice;

pub use crash::CrashGame;
pub use mines::MinesGame;
pub use plinko::PlinkoGame;
pub use dice::DiceGame;
