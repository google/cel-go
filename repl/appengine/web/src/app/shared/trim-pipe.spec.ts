import { TrimPipe } from './trim-pipe';

describe('TrimPipe', () => {
  it('create an instance', () => {
    const pipe = new TrimPipe();
    expect(pipe).toBeTruthy();
  });

  it('trims long string', () => {
    const pipe = new TrimPipe();
    expect(pipe.transform("test123456789", 7)).toBe("test...");
  });

  it('does not trim short string', () => {
    const pipe = new TrimPipe();
    expect(pipe.transform("test", 7)).toBe("test");
  });
});
