//! TigerCasino Security CLI
//! 
//! Command-line tool for security operations

use tigercasino_security::Security;

fn main() {
    println!("TigerCasino Security Module");
    println!("===========================");
    println!();
    
    // Example operations
    println!("1. Random Number Generation:");
    let random = Security::generate_random(1, 100);
    println!("   Random (1-100): {}", random);
    
    println!("\n2. UUID Generation:");
    let uuid = Security::generate_uuid();
    println!("   UUID: {}", uuid);
    
    println!("\n3. Password Hashing:");
    let password = "SecureP@ssw0rd123";
    match Security::hash_password(password) {
        Ok(hash) => {
            println!("   Hash: {}", &hash[..40]);
            match Security::verify_password(password, &hash) {
                Ok(valid) => println!("   Verified: {}", valid),
                Err(e) => println!("   Error: {}", e),
            }
        }
        Err(e) => println!("   Error: {}", e),
    }
    
    println!("\n4. Encryption/Decryption:");
    let key = [0u8; 32];
    let data = b"Sensitive casino data";
    match Security::encrypt(data, &key) {
        Ok(encrypted) => {
            println!("   Encrypted (base64): {}", Security::base64_encode(&encrypted));
            match Security::decrypt(&encrypted, &key) {
                Ok(decrypted) => println!("   Decrypted: {}", String::from_utf8_lossy(&decrypted)),
                Err(e) => println!("   Decryption error: {}", e),
            }
        }
        Err(e) => println!("   Encryption error: {}", e),
    }
    
    println!("\n5. HMAC Signature:");
    let message = b"Bet:100:User:123";
    let secret = b"casino_secret_key";
    let sig = Security::create_hmac(message, secret);
    println!("   Signature: {}", sig);
    println!("   Verified: {}", Security::verify_hmac(message, secret, &sig));
    
    println!("\n6. Fraud Detection:");
    let risk = Security::detect_fraud_pattern(50000.0, 10000.0, 0.98, 150);
    println!("   Risk Level: {}", risk.as_str());
    
    println!("\n7. SHA-256 Hash:");
    let data = b"Casino transaction data";
    println!("   Hash: {}", Security::sha256(data));
    
    println!("\nSecurity module ready for production use!");
}
