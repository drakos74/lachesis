// SPDX-License-Identifier: Unlicense OR MIT

package org.gioui;

import java.lang.Class;
import java.lang.IllegalAccessException;
import java.lang.InstantiationException;
import java.lang.ExceptionInInitializerError;
import java.lang.SecurityException;
import android.app.Activity;
import android.app.Fragment;
import android.app.FragmentManager;
import android.app.FragmentTransaction;
import android.content.Context;
import android.graphics.Rect;
import android.os.Build;
import android.os.Handler;
import android.text.Editable;
import android.util.AttributeSet;
import android.view.Choreographer;
import android.view.KeyCharacterMap;
import android.view.KeyEvent;
import android.view.MotionEvent;
import android.view.View;
import android.view.WindowInsets;
import android.view.Surface;
import android.view.SurfaceView;
import android.view.SurfaceHolder;
import android.view.inputmethod.BaseInputConnection;
import android.view.inputmethod.InputConnection;
import android.view.inputmethod.InputMethodManager;
import android.view.inputmethod.EditorInfo;

import java.io.UnsupportedEncodingException;

public class GioView extends SurfaceView implements Choreographer.FrameCallback {
	private final static Object initLock = new Object();
	private static boolean jniLoaded;

	private final SurfaceHolder.Callback callbacks;
	private final InputMethodManager imm;
	private final Handler handler;
	private long nhandle;

	private static synchronized void initialize(Context appCtx) {
		synchronized (initLock) {
			if (jniLoaded) {
				return;
			}
			String dataDir = appCtx.getFilesDir().getAbsolutePath();
			byte[] dataDirUTF8;
			try {
				dataDirUTF8 = dataDir.getBytes("UTF-8");
			} catch (UnsupportedEncodingException e) {
				throw new RuntimeException(e);
			}
			System.loadLibrary("gio");
			runGoMain(dataDirUTF8, appCtx);
			jniLoaded = true;
		}
	}

	public GioView(Context context) {
		this(context, null);
	}

	public GioView(Context context, AttributeSet attrs) {
		super(context, attrs);

		handler = new Handler();
		// Late initialization of the Go runtime to wait for a valid context.
		initialize(context.getApplicationContext());

		nhandle = onCreateView(this);
		imm = (InputMethodManager)context.getSystemService(Context.INPUT_METHOD_SERVICE);
		setFocusable(true);
		setFocusableInTouchMode(true);
		setOnFocusChangeListener(new View.OnFocusChangeListener() {
			@Override public void onFocusChange(View v, boolean focus) {
				GioView.this.onFocusChange(nhandle, focus);
			}
		});
		callbacks = new SurfaceHolder.Callback() {
			@Override public void surfaceCreated(SurfaceHolder holder) {
				// Ignore; surfaceChanged is guaranteed to be called immediately after this.
			}
			@Override public void surfaceChanged(SurfaceHolder holder, int format, int width, int height) {
				onSurfaceChanged(nhandle, getHolder().getSurface());
			}
			@Override public void surfaceDestroyed(SurfaceHolder holder) {
				onSurfaceDestroyed(nhandle);
			}
		};
		getHolder().addCallback(callbacks);
	}

	@Override public boolean onKeyDown(int keyCode, KeyEvent event) {
		onKeyEvent(nhandle, keyCode, event.getUnicodeChar(), event.getEventTime());
		return false;
	}

	@Override public boolean onTouchEvent(MotionEvent event) {
		// Ask for unbuffered events. Flutter and Chrome does it
		// so I assume its good for us as well.
		if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
			requestUnbufferedDispatch(event);
		}

		for (int j = 0; j < event.getHistorySize(); j++) {
			long time = event.getHistoricalEventTime(j);
			for (int i = 0; i < event.getPointerCount(); i++) {
				onTouchEvent(
						nhandle,
						event.ACTION_MOVE,
						event.getPointerId(i),
						event.getToolType(i),
						event.getHistoricalX(i, j),
						event.getHistoricalY(i, j),
						event.getButtonState(),
						time);
			}
		}
		int act = event.getActionMasked();
		int idx = event.getActionIndex();
		for (int i = 0; i < event.getPointerCount(); i++) {
			int pact = event.ACTION_MOVE;
			if (i == idx) {
				pact = act;
			}
			onTouchEvent(
					nhandle,
					act,
					event.getPointerId(i),
					event.getToolType(i),
					event.getX(i),
					event.getY(i),
					event.getButtonState(),
					event.getEventTime());
		}
		return true;
	}

	@Override public InputConnection onCreateInputConnection(EditorInfo outAttrs) {
		return new InputConnection(this);
	}

	void showTextInput() {
		post(new Runnable() {
			@Override public void run() {
				GioView.this.requestFocus();
				imm.showSoftInput(GioView.this, 0);
			}
		});
	}

	void hideTextInput() {
		post(new Runnable() {
			@Override public void run() {
				imm.hideSoftInputFromWindow(getWindowToken(), 0);
			}
		});
	}

	void postFrameCallbackOnMainThread() {
		handler.post(new Runnable() {
			@Override public void run() {
				postFrameCallback();
			}
		});
	}

	@Override protected boolean fitSystemWindows(Rect insets) {
		onWindowInsets(nhandle, insets.top, insets.right, insets.bottom, insets.left);
		return true;
	}

	void postFrameCallback() {
		Choreographer.getInstance().removeFrameCallback(this);
		Choreographer.getInstance().postFrameCallback(this);
	}

	@Override public void doFrame(long nanos) {
		onFrameCallback(nhandle, nanos);
	}

	int getDensity() {
		return getResources().getDisplayMetrics().densityDpi;
	}

	float getFontScale() {
		return getResources().getConfiguration().fontScale;
	}

	void start() {
		onStartView(nhandle);
	}

	void stop() {
		onStopView(nhandle);
	}

	void destroy() {
		getHolder().removeCallback(callbacks);
		onDestroyView(nhandle);
		nhandle = 0;
	}

	void configurationChanged() {
		onConfigurationChanged(nhandle);
	}

	void lowMemory() {
		onLowMemory();
	}

	boolean backPressed() {
		return onBack(nhandle);
	}

	public void registerFragment(String del) {
		final Class cls;
		try {
			cls = getContext().getClassLoader().loadClass(del);
		} catch (ClassNotFoundException e) {
			throw new RuntimeException("RegisterFragment: fragment class not found: " + e.getMessage());
		}

		handler.post(new Runnable() {
			public void run() {
				final Fragment frag;
				try {
					frag = (Fragment)cls.newInstance();
				} catch (IllegalAccessException | InstantiationException | ExceptionInInitializerError | SecurityException | ClassCastException e) {
					throw new RuntimeException("RegisterFragment: error instantiating fragment: " + e.getMessage());
				}
				final FragmentManager fm;
				try {
					fm = ((Activity)getContext()).getFragmentManager();
				} catch (ClassCastException e) {
					throw new RuntimeException("RegisterFragment: cannot get fragment manager from View Context: " + e.getMessage());
				}
				FragmentTransaction ft = fm.beginTransaction();
				ft.add(frag, del);
				ft.commitNow();
			}
		});
	}

	static private native long onCreateView(GioView view);
	static private native void onDestroyView(long handle);
	static private native void onStartView(long handle);
	static private native void onStopView(long handle);
	static private native void onSurfaceDestroyed(long handle);
	static private native void onSurfaceChanged(long handle, Surface surface);
	static private native void onConfigurationChanged(long handle);
	static private native void onWindowInsets(long handle, int top, int right, int bottom, int left);
	static private native void onLowMemory();
	static private native void onTouchEvent(long handle, int action, int pointerID, int tool, float x, float y, int buttons, long time);
	static private native void onKeyEvent(long handle, int code, int character, long time);
	static private native void onFrameCallback(long handle, long nanos);
	static private native boolean onBack(long handle);
	static private native void onFocusChange(long handle, boolean focus);
	static private native void runGoMain(byte[] dataDir, Context context);

	private static class InputConnection extends BaseInputConnection {
		private final Editable editable;

		InputConnection(View view) {
			// Passing false enables "dummy mode", where the BaseInputConnection
			// attempts to convert IME operations to key events.
			super(view, false);
			editable = Editable.Factory.getInstance().newEditable("");
		}

		@Override public Editable getEditable() {
			return editable;
		}
	}
}
